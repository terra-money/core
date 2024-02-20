package ante

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/terra-money/core/v2/x/smartaccount/types"
)

// SmartAccountAuthDecorator does authentication for smart accounts
type SmartAccountAuthDecorator struct {
	sak                       SmartAccountKeeper
	wk                        WasmKeeper
	ak                        authante.AccountKeeper
	defaultVerifySigDecorator sdk.AnteHandler
}

func NewSmartAccountAuthDecorator(
	sak SmartAccountKeeper,
	wk WasmKeeper,
	ak authante.AccountKeeper,
	sigGasConsumer func(meter sdk.GasMeter, sig signing.SignatureV2, params authtypes.Params) error,
	signModeHandler authsigning.SignModeHandler,
) SmartAccountAuthDecorator {
	if sigGasConsumer == nil {
		sigGasConsumer = authante.DefaultSigVerificationGasConsumer
	}
	defaultVerifySigDecorator := sdk.ChainAnteDecorators(
		authante.NewSetPubKeyDecorator(ak),
		authante.NewValidateSigCountDecorator(ak),
		authante.NewSigGasConsumeDecorator(ak, authante.DefaultSigVerificationGasConsumer),
		authante.NewSigVerificationDecorator(ak, signModeHandler),
	)
	return SmartAccountAuthDecorator{
		sak:                       sak,
		wk:                        wk,
		defaultVerifySigDecorator: defaultVerifySigDecorator,
	}
}

// AnteHandle checks if the tx provides sufficient fee to cover the required fee from the fee market.
func (sad SmartAccountAuthDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// check if the tx is from a smart account
	setting, err := sad.sak.GetSetting(ctx, tx.GetMsgs()[0].GetSigners()[0].String())
	if sdkerrors.ErrKeyNotFound.Is(err) {
		// run through the default handlers for signature verification
		newCtx, err := sad.defaultVerifySigDecorator(ctx, tx, simulate)
		if err != nil {
			return newCtx, err
		}
		// continue to the next handler after default signature verification
		return next(newCtx, tx, simulate)
	} else if err != nil {
		return ctx, err
	}

	if setting.Authorization != nil && len(setting.Authorization) > 0 {
		for _, auth := range setting.Authorization {
			// TODO: add code that calls authorization on contracts
			authMsg := types.Authorization{
				// TODO: check these fields
				Sender:  sad.ak.GetModuleAddress(types.ModuleName).String(),
				Account: tx.GetMsgs()[0].GetSigners()[0].String(),
				Data:    []byte(auth.InitMsg),
				// TODO: fill the below fields
				Signatures:  [][]byte{},
				SignedBytes: []byte{},
			}
			authMsgBz, err := json.Marshal(authMsg)
			if err != nil {
				return ctx, err
			}
			if _, err = sad.wk.Sudo(ctx, sdk.AccAddress(auth.ContractAddress), authMsgBz); err != nil {
				return ctx, err
			}
			if err != nil && setting.Fallback {
				return next(ctx, tx, simulate)
			} else if err != nil {
				return ctx, err
			}
		}
	}

	return next(ctx, tx, simulate)
}
