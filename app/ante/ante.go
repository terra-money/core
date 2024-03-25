package ante

import (
	ibcante "github.com/cosmos/ibc-go/v7/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	feesharekeeper "github.com/terra-money/core/v2/x/feeshare/keeper"

	smartaccountante "github.com/terra-money/core/v2/x/smartaccount/ante"
	smartaccountkeeper "github.com/terra-money/core/v2/x/smartaccount/keeper"

	"github.com/cosmos/cosmos-sdk/client"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	terrawasmkeeper "github.com/terra-money/core/v2/x/wasm/keeper"
)

// HandlerOptions extends the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	ante.HandlerOptions

	IBCkeeper          *ibckeeper.Keeper
	FeeShareKeeper     feesharekeeper.Keeper
	BankKeeper         bankKeeper.Keeper
	SmartAccountKeeper *smartaccountkeeper.Keeper
	WasmKeeper         *terrawasmkeeper.Keeper
	TxCounterStoreKey  storetypes.StoreKey
	WasmConfig         wasmTypes.WasmConfig
	TxConfig           client.TxConfig
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	if options.SmartAccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "smart account keeper is required for ante builder")
	}

	if options.WasmKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "wasm keeper is required for ante builder")
	}

	sigGasConsumer := options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit),
		wasmkeeper.NewCountTXDecorator(options.TxCounterStoreKey),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		// TODO: remove the following line after the migration to the new signature verification decorator is done
		// SetPubKeyDecorator must be called before all signature verification decorators
		// ante.NewSetPubKeyDecorator(options.AccountKeeper),
		// ante.NewValidateSigCountDecorator(options.AccountKeeper),
		// ante.NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		// ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		smartaccountante.NewSmartAccountAuthDecorator(*options.SmartAccountKeeper, options.WasmKeeper, options.AccountKeeper, sigGasConsumer, options.SignModeHandler),
		smartaccountante.NewPreTransactionHookDecorator(*options.SmartAccountKeeper, options.WasmKeeper),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCkeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
