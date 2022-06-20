package client

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"

	"github.com/terra-money/core/v2/app/ante/client/cli"
)

// ProposalHandler is the minimum commission update proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitMinimumCommissionUpdateProposal, emptyHandler)

func emptyHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "ante",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			rest.WriteErrorResponse(w, http.StatusNotFound, "end point not implemented")
			return
		},
	}
}
