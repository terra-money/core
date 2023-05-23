package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-money/core/v2/app/config"
)

func RegisterDenomsConfig() {
	// sdk.RegisterDenom(config.Luna, sdk.OneDec())
	// sdk.RegisterDenom(config.MilliLuna, sdk.NewDecWithPrec(1, 3))
	sdk.RegisterDenom(config.MicroLuna, sdk.NewDecWithPrec(1, 6))
	// sdk.RegisterDenom(config.NanoLuna, sdk.NewDecWithPrec(1, 9))
}
