package ante

import (
	sdkerrors "cosmossdk.io/errors"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrortypes "github.com/cosmos/cosmos-sdk/types/errors"
	tx2 "github.com/cosmos/cosmos-sdk/types/tx"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	smartaccounttypes "github.com/terra-money/core/v2/x/smartaccount/types"
)

// SmartAccountCheckDecorator does authentication for smart accounts
type PreTransactionHookDecorator struct {
	smartaccountKeeper SmartAccountKeeper
	wasmKeeper         WasmKeeper
}

func NewPreTransactionHookDecorator(sak SmartAccountKeeper, wk WasmKeeper) PreTransactionHookDecorator {
	return PreTransactionHookDecorator{
		smartaccountKeeper: sak,
		wasmKeeper:         wk,
	}
}

// AnteHandle checks if the tx provides sufficient fee to cover the required fee from the fee market.
func (pth PreTransactionHookDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	setting, ok := ctx.Value(smartaccounttypes.ModuleName).(smartaccounttypes.Setting)
	if !ok {
		return next(ctx, tx, simulate)
	}

	if setting.PreTransaction != nil && len(setting.PreTransaction) > 0 {
		for _, preTx := range setting.PreTransaction {
			contractAddr, err := sdk.AccAddressFromBech32(preTx)
			if err != nil {
				return ctx, err
			}
			data, err := BuildPreTransactionHookMsg(tx)
			_, err = pth.wasmKeeper.Sudo(ctx, contractAddr, data)
			if err != nil {
				return ctx, err
			}
		}
	}

	return next(ctx, tx, simulate)
}

// TODO: to refactor
func BuildPreTransactionHookMsg(tx sdk.Tx) ([]byte, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return nil, sdkerrors.Wrap(sdkerrortypes.ErrInvalidType, "expected SigVerifiableTx")
	}

	// Signer here is the account that the state transition is affecting
	// e.g. Account that is transferring some Coins
	signers := sigTx.GetSigners()
	// Current only supports one signer (TODO review in the future)
	if len(signers) != 1 {
		return nil, sdkerrors.Wrap(sdkerrortypes.ErrorInvalidSigner, "only one signer is supported")
	}

	// Sender here is the account that signed the transaction
	// Could be different from the account above (confusingly named signer)
	signatures, _ := sigTx.GetSignaturesV2()
	if len(signatures) == 0 {
		return nil, sdkerrors.Wrap(sdkerrortypes.ErrNoSignatures, "no signatures found")
	}
	senderAddr, err := sdk.AccAddressFromHexUnsafe(signatures[0].PubKey.Address().String())
	if err != nil {
		return nil, err
	}

	msgs := sigTx.GetMsgs()
	anyMsgs, err := tx2.SetMsgs(msgs)
	if err != nil {
		return nil, err
	}
	var stargateMsgs []wasmvmtypes.CosmosMsg
	for _, msg := range anyMsgs {
		stargateMsg := wasmvmtypes.StargateMsg{
			TypeURL: msg.TypeUrl,
			Value:   msg.Value,
		}
		stargateMsgs = append(stargateMsgs, wasmvmtypes.CosmosMsg{
			Stargate: &stargateMsg,
		})
	}
	preTx := smartaccounttypes.PreTransaction{
		Sender:   senderAddr.String(),
		Account:  signers[0].String(),
		Messages: stargateMsgs,
	}
	msg := smartaccounttypes.SudoMsg{PreTransaction: &preTx}
	return json.Marshal(msg)
}
