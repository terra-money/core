package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeMinimumCommissionUpdate string = "MinimumCommissionUpdate"
)

func NewMinimumCommissionUpdateProposal(title, description string, minimumCommission sdk.Dec) govtypes.Content {
	return &MinimumCommissionUpdateProposal{title, description, minimumCommission}
}

// Implements Proposal Interface
var _ govtypes.Content = &MinimumCommissionUpdateProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeMinimumCommissionUpdate)
	govtypes.RegisterProposalTypeCodec(&MinimumCommissionUpdateProposal{}, "ante/MinimumCommissionUpdateProposal")
}

func (sup *MinimumCommissionUpdateProposal) GetTitle() string       { return sup.Title }
func (sup *MinimumCommissionUpdateProposal) GetDescription() string { return sup.Description }
func (sup *MinimumCommissionUpdateProposal) ProposalRoute() string  { return RouterKey }
func (sup *MinimumCommissionUpdateProposal) ProposalType() string {
	return ProposalTypeMinimumCommissionUpdate
}
func (miup *MinimumCommissionUpdateProposal) ValidateBasic() error {

	if err := validateMinimumCommission(miup.MinimumCommission); err != nil {
		return err
	}
	return govtypes.ValidateAbstract(miup)
}

func (sup MinimumCommissionUpdateProposal) String() string {
	return fmt.Sprintf(`Minimum Commission Update Proposal:
  Title:       %s
  Description: %s
`, sup.Title, sup.Description)
}
