package ante

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	antekeeper "github.com/terra-money/core/v2/app/ante/keeper"
)

// MinCommissionDecorator checks whether the validator's commission rate
// is smaller than hard limit(= MinimumCommissionRate) or not
type MinCommissionDecorator struct {
	cdc        codec.BinaryCodec
	anteKeeper *antekeeper.Keeper
}

// NewMinCommissionDecorator return MinCommissionDecorator instance
func NewMinCommissionDecorator(cdc codec.BinaryCodec, anteKeeper *antekeeper.Keeper) MinCommissionDecorator {
	return MinCommissionDecorator{cdc, anteKeeper}
}

// AnteHandle implement interface
func (min MinCommissionDecorator) AnteHandle(
	ctx sdk.Context, tx sdk.Tx,
	simulate bool, next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	if min.anteKeeper != nil && min.anteKeeper.MinimumCommissionEnforced(ctx) {
		minimumCommission := min.anteKeeper.GetMinimumCommission(ctx)

		msgs := tx.GetMsgs()
		validMsg := func(m sdk.Msg) error {
			switch msg := m.(type) {
			case *stakingtypes.MsgCreateValidator:
				// prevent new validators joining the set with
				// commission set below MinimumCommission
				c := msg.Commission
				if c.Rate.LT(minimumCommission) {
					return sdkerrors.Wrap(sdkerrors.ErrUnauthorized,
						fmt.Sprintf("commission can't be lower than %s%%", minimumCommission.MulInt64(100).String()),
					)
				}
			case *stakingtypes.MsgEditValidator:
				// if commission rate is nil, it means only
				// other fields are affected - skip
				if msg.CommissionRate == nil {
					break
				}

				if msg.CommissionRate.LT(minimumCommission) {
					return sdkerrors.Wrap(sdkerrors.ErrUnauthorized,
						fmt.Sprintf("commission can't be lower than %s%%", minimumCommission.String()),
					)
				}
			}

			return nil
		}

		validAuthz := func(execMsg *authz.MsgExec) error {
			for _, v := range execMsg.Msgs {
				var innerMsg sdk.Msg
				err := min.cdc.UnpackAny(v, &innerMsg)
				if err != nil {
					return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "cannot unmarshal authz exec msgs")
				}

				err = validMsg(innerMsg)
				if err != nil {
					return err
				}
			}

			return nil
		}

		for _, m := range msgs {
			if msg, ok := m.(*authz.MsgExec); ok {
				if err := validAuthz(msg); err != nil {
					return ctx, err
				}
				continue
			}

			// validate normal msgs
			err = validMsg(m)
			if err != nil {
				return ctx, err
			}
		}
	}

	return next(ctx, tx, simulate)
}
