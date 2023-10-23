package params

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/terra-money/core/v2/app/config"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterAddressesConfig() *sdk.Config {
	sdkConfig := sdk.GetConfig()
	sdkConfig.SetCoinType(config.CoinType)

	accountPubKeyPrefix := config.AccountAddressPrefix + "pub"
	validatorAddressPrefix := config.AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := config.AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := config.AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := config.AccountAddressPrefix + "valconspub"

	sdkConfig.SetBech32PrefixForAccount(config.AccountAddressPrefix, accountPubKeyPrefix)
	sdkConfig.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	sdkConfig.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	sdkConfig.SetAddressVerifier(wasmtypes.VerifyAddressLen())

	return sdkConfig
}
