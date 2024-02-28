package ante

import (
	"encoding/json"
	"fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
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
	smartAccountKeeper        SmartAccountKeeper
	wasmKeeper                WasmKeeper
	accountKeeper             authante.AccountKeeper
	signModeHandler           authsigning.SignModeHandler
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
		authante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		authante.NewSigVerificationDecorator(ak, signModeHandler),
	)
	return SmartAccountAuthDecorator{
		smartAccountKeeper:        sak,
		wasmKeeper:                wk,
		accountKeeper:             ak,
		signModeHandler:           signModeHandler,
		defaultVerifySigDecorator: defaultVerifySigDecorator,
	}
}

// AnteHandle checks if the tx provides sufficient fee to cover the required fee from the fee market.
func (sad SmartAccountAuthDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// skip the smartaccount auth decorator if the tx is a simulation
	if simulate {
		return next(ctx, tx, simulate)
	}

	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.ErrInvalidType.Wrap("expected SigVerifiableTx")
	}

	// Signer here is the account that the state transition is affecting
	// e.g. Account that is transferring some Coins
	signers := sigTx.GetSigners()
	account := signers[0].String()

	// check if the tx is from a smart account
	setting, err := sad.smartAccountKeeper.GetSetting(ctx, account)
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

	// Current smartaccount only supports one signer (TODO review in the future)
	if len(signers) != 1 {
		return ctx, sdkerrors.ErrorInvalidSigner.Wrap("only one account is supported (sigTx.GetSigners()!= 1)")
	}

	// Sender here is the account that signed the transaction
	// Could be different from the account above (confusingly named signer)
	signatures, err := sigTx.GetSignaturesV2()
	if err != nil {
		return ctx, err
	}
	if len(signatures) == 0 {
		return ctx, sdkerrors.ErrNoSignatures.Wrap("no signatures found")
	} else if len(signatures) > 1 {
		// TODO: remove when support multi sig auth
		return ctx, sdkerrors.ErrUnauthorized.Wrap("multiple signatures not supported")
	}

	signature := signatures[0]
	signaturesBs := [][]byte{}

	senderAddr, err := sdk.AccAddressFromHexUnsafe(signature.PubKey.Address().String())
	if err != nil {
		return ctx, err
	}

	acc, err := authante.GetSignerAcc(ctx, sad.accountKeeper, senderAddr)
	if err != nil {
		return ctx, err
	}
	var accNum uint64
	if ctx.BlockHeight() != 0 {
		accNum = acc.GetAccountNumber()
	}

	signerData := authsigning.SignerData{
		Address:       senderAddr.String(),
		ChainID:       ctx.ChainID(),
		AccountNumber: accNum,
		Sequence:      acc.GetSequence(),
		PubKey:        signature.PubKey,
	}

	signatureBz, err := signatureDataToBz(signature.Data)
	if err != nil {
		return ctx, err
	}
	signedBytes, err := GetSignBytesArr(signature.PubKey, signerData, signature.Data, sad.signModeHandler, tx)
	if err != nil {
		return ctx, err
	}
	signaturesBs = append(signaturesBs, signatureBz...)

	ctx = ctx.WithValue(types.ModuleName, setting)
	if setting.Authorization != nil && len(setting.Authorization) > 0 {
		success := false
		for _, auth := range setting.Authorization {
			authMsg := types.Authorization{
				Sender:      senderAddr.String(),
				Account:     account,
				Data:        []byte(auth.InitMsg),
				Signatures:  signaturesBs,
				SignedBytes: signedBytes,
			}
			authMsgBz, err := json.Marshal(authMsg)
			if err != nil {
				return ctx, err
			}
			contractAddr, err := sdk.AccAddressFromBech32(auth.ContractAddress)
			if err != nil {
				return ctx, err
			}
			_, err = sad.wasmKeeper.Sudo(ctx, contractAddr, authMsgBz)
			// so long as one of the authorization is successful, we're good
			if err == nil {
				success = true
				break
			}
			if err != nil && setting.Fallback {
				return next(ctx, tx, simulate)
			} else if err != nil {
				return ctx, err
			}
		}
		if success {
			return next(ctx, tx, simulate)
		} else if !setting.Fallback {
			return ctx, sdkerrors.ErrUnauthorized.Wrap("authorization failed")
		}
	}
	// run through the default handlers for signature verification
	newCtx, err := sad.defaultVerifySigDecorator(ctx, tx, simulate)
	if err != nil {
		return newCtx, err
	}
	// continue to the next handler after default signature verification
	return next(newCtx, tx, simulate)
}

// signatureDataToBz converts a SignatureData into raw bytes signature.
// For SingleSignatureData, it returns the signature raw bytes.
// For MultiSignatureData, it returns an array of all individual signatures,
// as well as the aggregated signature.
func signatureDataToBz(data signing.SignatureData) ([][]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("got empty SignatureData")
	}

	switch data := data.(type) {
	case *signing.SingleSignatureData:
		return [][]byte{data.Signature}, nil
	case *signing.MultiSignatureData:
		sigs := [][]byte{}
		var err error

		for _, d := range data.Signatures {
			nestedSigs, err := signatureDataToBz(d)
			if err != nil {
				return nil, err
			}
			sigs = append(sigs, nestedSigs...)
		}

		multisig := cryptotypes.MultiSignature{
			Signatures: sigs,
		}
		aggregatedSig, err := multisig.Marshal()
		if err != nil {
			return nil, err
		}
		sigs = append(sigs, aggregatedSig)

		return sigs, nil
	default:
		return nil, sdkerrors.ErrInvalidType.Wrapf("unexpected signature data type %T", data)
	}
}

func GetSignBytesArr(pubKey cryptotypes.PubKey, signerData authsigning.SignerData, sigData signing.SignatureData, handler authsigning.SignModeHandler, tx sdk.Tx) (signersBytes [][]byte, err error) {
	switch data := sigData.(type) {
	case *signing.SingleSignatureData:
		signBytes, err := handler.GetSignBytes(data.SignMode, signerData, tx)
		if err != nil {
			return nil, err
		}
		if !pubKey.VerifySignature(signBytes, data.Signature) {
			return nil, fmt.Errorf("unable to verify single signer signature")
		}
		return [][]byte{signBytes}, nil

	case *signing.MultiSignatureData:
		multiPK, ok := pubKey.(multisig.PubKey)
		if !ok {
			return nil, fmt.Errorf("expected %T, got %T", (multisig.PubKey)(nil), pubKey)
		}
		return GetMultiSigSignBytes(multiPK, data, signerData, handler, tx)
	default:
		return nil, fmt.Errorf("unexpected SignatureData %T", sigData)
	}
}

func GetMultiSigSignBytes(multiPK multisig.PubKey, sig *signing.MultiSignatureData, signerData authsigning.SignerData, handler authsigning.SignModeHandler, tx sdk.Tx) (signersBytes [][]byte, err error) {
	bitarray := sig.BitArray
	sigs := sig.Signatures
	size := bitarray.Count()
	pubKeys := multiPK.GetPubKeys()
	// ensure bit array is the correct size
	if len(pubKeys) != size {
		return nil, fmt.Errorf("bit array size is incorrect, expecting: %d", len(pubKeys))
	}
	// ensure size of signature list
	if len(sigs) < int(multiPK.GetThreshold()) || len(sigs) > size {
		return nil, fmt.Errorf("signature size is incorrect %d", len(sigs))
	}
	// ensure at least k signatures are set
	if bitarray.NumTrueBitsBefore(size) < int(multiPK.GetThreshold()) {
		return nil, fmt.Errorf("not enough signatures set, have %d, expected %d", bitarray.NumTrueBitsBefore(size), int(multiPK.GetThreshold()))
	}
	// index in the list of signatures which we are concerned with.
	sigIndex := 0
	for i := 0; i < size; i++ {
		if bitarray.GetIndex(i) {
			si := sig.Signatures[sigIndex]
			switch si := si.(type) {
			case *signing.SingleSignatureData:
				signerBytes, err := handler.GetSignBytes(si.SignMode, signerData, tx)
				if err != nil {
					return nil, err
				}
				signersBytes = append(signersBytes, signerBytes)
			case *signing.MultiSignatureData:
				nestedMultisigPk, ok := pubKeys[i].(multisig.PubKey)
				if !ok {
					return nil, fmt.Errorf("unable to parse pubkey of index %d", i)
				}
				signersBytesHold, err := GetMultiSigSignBytes(nestedMultisigPk, si, signerData, handler, tx)
				if err != nil {
					return nil, err
				}
				signersBytes = append(signersBytes, signersBytesHold...)
			default:
				return nil, fmt.Errorf("improper signature data type for index %d", sigIndex)
			}
			sigIndex++
		}
	}
	return signersBytes, nil
}
