package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/terra-money/core/v2/app/ante/keeper"
	"github.com/terra-money/core/v2/app/ante/types"
)

// NewMinimumCommissionUpdateProposalHandler creates a governance handler to manage new proposal types.
// It updates all validators minimum commission to be bigger than the given value and updates ante stored
// minimum commission value to this value to block all messages which are making validator's commission
// smaller than the value.
func NewMinimumCommissionUpdateProposalHandler(k keeper.Keeper, sk stakingkeeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.MinimumCommissionUpdateProposal:
			return handleMinimumCommissionUpdateProposal(ctx, k, sk, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ante proposal content type: %T", c)
		}
	}
}

func handleMinimumCommissionUpdateProposal(
	ctx sdk.Context,
	k keeper.Keeper,
	sk stakingkeeper.Keeper,
	p *types.MinimumCommissionUpdateProposal) error {
	// update all validators minimum commission
	minimumCommission := p.MinimumCommission

	allValidators := sk.GetAllValidators(ctx)
	for _, validator := range allValidators {
		// increase commission rate
		if validator.Commission.CommissionRates.Rate.LT(minimumCommission) {

			// call the before-modification hook since we're about to update the commission
			sk.BeforeValidatorModified(ctx, validator.GetOperator())

			validator.Commission.Rate = minimumCommission
			validator.Commission.UpdateTime = ctx.BlockHeader().Time
		}

		// increase max commission rate
		if validator.Commission.CommissionRates.MaxRate.LT(minimumCommission) {
			validator.Commission.MaxRate = minimumCommission
		}

		sk.SetValidator(ctx, validator)
	}

	// update minimum commission
	k.SetMinimumCommission(ctx, p.MinimumCommission)
	return nil
}
