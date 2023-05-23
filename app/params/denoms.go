package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/v2/app/config"
)

func RegisterDenomsConfig() error {
	// sdk.RegisterDenom(config.Luna, sdk.OneDec())
	// sdk.RegisterDenom(config.MilliLuna, sdk.NewDecWithPrec(1, 3))
	err := sdk.RegisterDenom(config.MicroLuna, sdk.NewDecWithPrec(1, 6))
	if err != nil {
		return err
	}
	// sdk.RegisterDenom(config.NanoLuna, sdk.NewDecWithPrec(1, 9))

	return nil
}
