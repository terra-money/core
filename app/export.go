package app

import (
	"encoding/json"
	"log"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ExportAppStateAndValidators exports the state of the application for a genesis
// file.
func (app *TerraApp) ExportAppStateAndValidators(
	forZeroHeight bool,
	jailAllowedAddrs []string,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	// We export at last height + 1, because that's the height at which
	// Tendermint will start InitChain.
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs)
	}

	genState := app.mm.ExportGenesis(ctx, app.appCodec)
	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(ctx, app.Keepers.StakingKeeper)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}
	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, nil
}

// prepare for fresh start at zero height
// NOTE zero height genesis is a temporary feature which will be deprecated
//
//	in favour of export at a block height
func (app *TerraApp) prepForZeroHeightGenesis(ctx sdk.Context, jailAllowedAddrs []string) {
	applyAllowedAddrs := false

	// check if there is a allowed address list
	if len(jailAllowedAddrs) > 0 {
		applyAllowedAddrs = true
	}

	allowedAddrsMap := make(map[string]bool)

	for _, addr := range jailAllowedAddrs {
		_, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			log.Fatal(err)
		}
		allowedAddrsMap[addr] = true
	}

	/* Just to be safe, assert the invariants on current state. */
	app.Keepers.CrisisKeeper.AssertInvariants(ctx)

	/* Handle fee distribution state. */

	// withdraw all validator commission
	app.Keepers.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		_, err := app.Keepers.DistrKeeper.WithdrawValidatorCommission(ctx, val.GetOperator())
		if err != nil {
			panic(err)
		}
		return false
	})

	// withdraw all delegator rewards
	dels := app.Keepers.StakingKeeper.GetAllDelegations(ctx)
	for _, delegation := range dels {
		_, err := app.Keepers.DistrKeeper.WithdrawDelegationRewards(ctx, delegation.GetDelegatorAddr(), delegation.GetValidatorAddr())
		if err != nil {
			panic(err)
		}
	}

	// clear validator slash events
	app.Keepers.DistrKeeper.DeleteAllValidatorSlashEvents(ctx)

	// clear validator historical rewards
	app.Keepers.DistrKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	// set context height to zero
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// reinitialize all validators
	app.Keepers.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		// donate any unwithdrawn outstanding reward fraction tokens to the community pool
		scraps := app.Keepers.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, val.GetOperator())
		feePool := app.Keepers.DistrKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		app.Keepers.DistrKeeper.SetFeePool(ctx, feePool)

		err := app.Keepers.DistrKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		if err != nil {
			panic(err)
		}
		return false
	})

	// reinitialize all delegations
	for _, del := range dels {
		err := app.Keepers.DistrKeeper.Hooks().BeforeDelegationCreated(ctx, del.GetDelegatorAddr(), del.GetValidatorAddr())
		if err != nil {
			panic(err)
		}
		err = app.Keepers.DistrKeeper.Hooks().AfterDelegationModified(ctx, del.GetDelegatorAddr(), del.GetValidatorAddr())
		if err != nil {
			panic(err)
		}
	}

	// reset context height
	ctx = ctx.WithBlockHeight(height)

	/* Handle staking state. */

	// iterate through redelegations, reset creation height
	app.Keepers.StakingKeeper.IterateRedelegations(ctx, func(_ int64, red stakingtypes.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		app.Keepers.StakingKeeper.SetRedelegation(ctx, red)
		return false
	})

	// iterate through unbonding delegations, reset creation height
	app.Keepers.StakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd stakingtypes.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		app.Keepers.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
		return false
	})

	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := ctx.KVStore(app.keys[stakingtypes.StoreKey])
	iter := sdk.KVStoreReversePrefixIterator(store, stakingtypes.ValidatorsKey)
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(iter.Key()[1:])
		validator, found := app.Keepers.StakingKeeper.GetValidator(ctx, addr)
		if !found {
			panic("expected validator, not found")
		}

		validator.UnbondingHeight = 0
		if applyAllowedAddrs && !allowedAddrsMap[addr.String()] {
			validator.Jailed = true
		}

		app.Keepers.StakingKeeper.SetValidator(ctx, validator)
		counter++
	}

	err := iter.Close()
	if err != nil {
		panic(err)
	}

	if _, err := app.Keepers.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx); err != nil {
		panic(err)
	}

	/* Handle slashing state. */

	// reset start height on signing infos
	app.Keepers.SlashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			app.Keepers.SlashingKeeper.SetValidatorSigningInfo(ctx, addr, info)
			return false
		},
	)
}
