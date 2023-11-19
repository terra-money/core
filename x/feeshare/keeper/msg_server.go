package keeper

import (
	"context"

	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/terra-money/core/v2/x/feeshare/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) GetIfContractWasCreatedFromFactory(ctx sdk.Context, msgSender sdk.AccAddress, info *wasmTypes.ContractInfo) bool {
	// Gov Module Admin
	govMod := k.accountKeeper.GetModuleAddress(govtypes.ModuleName).String()
	if info.Admin == govMod {
		// only register to self.
		return true
	}

	if len(info.Admin) == 0 {
		// There is no admin. Return if the creator is a contract or normal user.
		creator, err := sdk.AccAddressFromBech32(info.Creator)
		if err != nil {
			return false
		}

		// is factory if creator is a contract.
		return k.wasmKeeper.HasContractInfo(ctx, creator)
	}

	// There is an admin
	admin, err := sdk.AccAddressFromBech32(info.Admin)
	if err != nil {
		return false
	}

	if admin.String() == msgSender.String() {
		return false
	}

	return k.wasmKeeper.HasContractInfo(ctx, admin)
}

// GetContractAdminOrCreatorAddress ensures the deployer is the contract's admin OR creator if no admin is set for all msg_server feeshare functions.
func (k Keeper) GetContractAdminOrCreatorAddress(ctx sdk.Context, contract sdk.AccAddress, deployer string) (sdk.AccAddress, error) {
	// Ensure the deployer address is valid
	_, err := sdk.AccAddressFromBech32(deployer)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid deployer address %s", deployer)
	}

	// Retrieve contract info
	info := k.wasmKeeper.GetContractInfo(ctx, contract)

	// Check if the contract has an admin
	if len(info.Admin) == 0 {
		// No admin, so check if the deployer is the creator of the contract
		if info.Creator != deployer {
			return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "you are not the creator of this contract %s", info.Creator)
		}

		_, err := sdk.AccAddressFromBech32(info.Creator)
		if err != nil {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address %s", info.Creator)
		}

		// Deployer is the creator, return the controlling account as the creator's address
		return sdk.AccAddressFromBech32(info.Creator)
	}

	// Admin is set, so check if the deployer is the admin of the contract
	if info.Admin != deployer {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "you are not an admin of this contract %s", deployer)
	}

	// Verify the admin address is valid
	_, err = sdk.AccAddressFromBech32(info.Admin)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid admin address %s", info.Admin)
	}

	return sdk.AccAddressFromBech32(info.Admin)
}

// RegisterFeeShare registers a contract to receive transaction fees
func (k Keeper) RegisterFeeShare(
	goCtx context.Context,
	msg *types.MsgRegisterFeeShare,
) (*types.MsgRegisterFeeShareResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if !params.EnableFeeShare {
		return nil, types.ErrFeeShareDisabled
	}

	// Get Contract
	contract, err := sdk.AccAddressFromBech32(msg.ContractAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", err)
	}

	// Check if contract is already registered
	if k.IsFeeShareRegistered(ctx, contract) {
		return nil, errorsmod.Wrapf(types.ErrFeeShareAlreadyRegistered, "contract is already registered %s", contract)
	}

	// Get the withdraw address of the contract
	withdrawer, err := sdk.AccAddressFromBech32(msg.WithdrawerAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid withdrawer address %s", msg.WithdrawerAddress)
	}

	// ensure msg.DeployerAddress is  valid
	msgSender, err := sdk.AccAddressFromBech32(msg.DeployerAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid deployer address %s", msg.DeployerAddress)
	}

	var deployer sdk.AccAddress

	if k.GetIfContractWasCreatedFromFactory(ctx, msgSender, k.wasmKeeper.GetContractInfo(ctx, contract)) {
		// Anyone is allowed to register a contract to itself if it was created from a factory contract
		if msg.WithdrawerAddress != msg.ContractAddress {
			return nil, errorsmod.Wrapf(types.ErrFeeShareInvalidWithdrawer, "withdrawer address must be the same as the contract address if it is from a factory contract withdraw:%s contract:%s", msg.WithdrawerAddress, msg.ContractAddress)
		}

		// set the deployer address to the contract address so it can self register
		msg.DeployerAddress = msg.ContractAddress
		deployer, err = sdk.AccAddressFromBech32(msg.DeployerAddress)
		if err != nil {
			return nil, err
		}
	} else {
		// Check that the person who signed the message is the wasm contract admin or creator (if no admin)
		deployer, err = k.GetContractAdminOrCreatorAddress(ctx, contract, msg.DeployerAddress)
		if err != nil {
			return nil, err
		}
	}

	// prevent storing the same address for deployer and withdrawer
	feeshare := types.NewFeeShare(contract, deployer, withdrawer)
	k.SetFeeShare(ctx, feeshare)
	k.SetDeployerMap(ctx, deployer, contract)
	k.SetWithdrawerMap(ctx, withdrawer, contract)

	k.Logger(ctx).Debug(
		"registering contract for transaction fees",
		"contract", msg.ContractAddress,
		"deployer", msg.DeployerAddress,
		"withdraw", msg.WithdrawerAddress,
	)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeRegisterFeeShare,
				// sdk.NewAttribute(sdk.AttributeKeySender, msg.DeployerAddress), // SDK v47
				sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress),
				sdk.NewAttribute(types.AttributeKeyWithdrawerAddress, msg.WithdrawerAddress),
			),
		},
	)

	return &types.MsgRegisterFeeShareResponse{}, nil
}

// UpdateFeeShare updates the withdraw address of a given FeeShare. If the given
// withdraw address is empty or the same as the deployer address, the withdraw
// address is removed.
func (k Keeper) UpdateFeeShare(
	goCtx context.Context,
	msg *types.MsgUpdateFeeShare,
) (*types.MsgUpdateFeeShareResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if !params.EnableFeeShare {
		return nil, types.ErrFeeShareDisabled
	}

	contract, err := sdk.AccAddressFromBech32(msg.ContractAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid contract address (%s)", err,
		)
	}

	feeshare, found := k.GetFeeShare(ctx, contract)
	if !found {
		return nil, errorsmod.Wrapf(
			types.ErrFeeShareContractNotRegistered,
			"contract %s is not registered", msg.ContractAddress,
		)
	}

	// feeshare with the given withdraw address is already registered
	if msg.WithdrawerAddress == feeshare.WithdrawerAddress {
		return nil, errorsmod.Wrapf(types.ErrFeeShareAlreadyRegistered, "feeshare with withdraw address %s is already registered", msg.WithdrawerAddress)
	}

	// Check that the person who signed the message is the wasm contract admin, if so return the deployer address
	_, err = k.GetContractAdminOrCreatorAddress(ctx, contract, msg.DeployerAddress)
	if err != nil {
		return nil, err
	}

	withdrawAddr, err := sdk.AccAddressFromBech32(feeshare.WithdrawerAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid withdrawer address (%s)", err,
		)
	}
	newWithdrawAddr, err := sdk.AccAddressFromBech32(msg.WithdrawerAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid WithdrawerAddress %s", msg.WithdrawerAddress)
	}

	k.DeleteWithdrawerMap(ctx, withdrawAddr, contract)
	k.SetWithdrawerMap(ctx, newWithdrawAddr, contract)

	// update feeshare
	feeshare.WithdrawerAddress = newWithdrawAddr.String()
	k.SetFeeShare(ctx, feeshare)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeUpdateFeeShare,
				// sdk.NewAttribute(sdk.AttributeKeySender, msg.DeployerAddress), // SDK v47
				sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress),
				sdk.NewAttribute(types.AttributeKeyWithdrawerAddress, msg.WithdrawerAddress),
			),
		},
	)

	return &types.MsgUpdateFeeShareResponse{}, nil
}

// CancelFeeShare deletes the FeeShare for a given contract
func (k Keeper) CancelFeeShare(
	goCtx context.Context,
	msg *types.MsgCancelFeeShare,
) (*types.MsgCancelFeeShareResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if !params.EnableFeeShare {
		return nil, types.ErrFeeShareDisabled
	}

	contract, err := sdk.AccAddressFromBech32(msg.ContractAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", err)
	}

	fee, found := k.GetFeeShare(ctx, contract)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrFeeShareContractNotRegistered, "contract %s is not registered", msg.ContractAddress)
	}

	// Check that the person who signed the message is the wasm contract admin, if so return the deployer address
	_, err = k.GetContractAdminOrCreatorAddress(ctx, contract, msg.DeployerAddress)
	if err != nil {
		return nil, err
	}

	k.DeleteFeeShare(ctx, fee)
	k.DeleteDeployerMap(
		ctx,
		fee.GetDeployerAddr(),
		contract,
	)

	withdrawAddr := fee.GetWithdrawerAddr()
	if withdrawAddr != nil {
		k.DeleteWithdrawerMap(
			ctx,
			withdrawAddr,
			contract,
		)
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeCancelFeeShare,
				// sdk.NewAttribute(sdk.AttributeKeySender, msg.DeployerAddress), // SDK v47
				sdk.NewAttribute(types.AttributeKeyContract, msg.ContractAddress),
			),
		},
	)

	return &types.MsgCancelFeeShareResponse{}, nil
}

func (k Keeper) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
