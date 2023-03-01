package rpc

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"net/http"
)

type HealthcheckResponse struct {
	Health string `json:"health"`
}

func RegisterHealthcheckRoute(context client.Context, r *mux.Router) {
	r.HandleFunc("/health", NodeHealthRequestHandlerFn(context)).Methods("GET")
}

// NodeHealthRequestHandlerFn
// REST handler for node health check - aws recognizes only http status codes
func NodeHealthRequestHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := getNodeStatus(clientCtx)
		if CheckInternalServerError(w, err) {
			return
		}
		if status.SyncInfo.CatchingUp {
			WriteErrorResponse(w, http.StatusServiceUnavailable, "NOK")
		} else {
			PostProcessResponseBare(w, clientCtx, HealthcheckResponse{Health: "OK"})
		}
	}
}

func getNodeStatus(clientCtx client.Context) (*ctypes.ResultStatus, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return &ctypes.ResultStatus{}, err
	}

	return node.Status(context.Background())
}
