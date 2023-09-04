package bank

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank/exported"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	customalliancemod "github.com/terra-money/alliance/custom/bank"
	custombankkeeper "github.com/terra-money/core/v2/custom/bank/keeper"
)

// AppModule wraps around the bank module and the bank keeper to return the right total supply ignoring bonded tokens
// that the alliance module minted to rebalance the voting power
// It modifies the TotalSupply and SupplyOf GRPC queries
type AppModule struct {
	customalliancemod.AppModule
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper custombankkeeper.Keeper, accountKeeper types.AccountKeeper, ss exported.Subspace) AppModule {
	mod := customalliancemod.NewAppModule(cdc, keeper.Keeper, accountKeeper, ss)

	return AppModule{
		AppModule: mod,
	}
}
