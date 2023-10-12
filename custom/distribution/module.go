package distribution

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank/exported"

	distributionmod "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	distrkeeper "github.com/terra-money/core/v2/custom/distribution/keeper"
)

// AppModule wraps around the bank module and the bank keeper to return the right total supply ignoring bonded tokens
// that the alliance module minted to rebalance the voting power
// It modifies the TotalSupply and SupplyOf GRPC queries
type AppModule struct {
	distributionmod.AppModule
	keeper distrkeeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec,
	keeper distrkeeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	ss exported.Subspace,
) AppModule {
	appModule := distributionmod.NewAppModule(
		cdc,
		keeper.Keeper,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		ss,
	)

	return AppModule{
		AppModule: appModule,
		keeper:    keeper,
	}
}
