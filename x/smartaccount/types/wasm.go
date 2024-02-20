package types

import "github.com/CosmWasm/wasmvm/types"

type Authorization struct {
	Sender      string   `json:"sender"`
	Account     string   `json:"account"`
	Signatures  [][]byte `json:"signatures"`
	SignedBytes []byte   `json:"signed_bytes"`
	Data        []byte   `json:"data"`
}

type PreTransaction struct {
	Sender   string              `json:"sender"`
	Account  string              `json:"account"`
	Messages []types.StargateMsg `json:"messages"`
}

type PostTransaction struct {
	Sender  string              `json:"sender"`
	Account string              `json:"account"`
	Msgs    []types.StargateMsg `json:"msgs"`
}
