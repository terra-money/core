package types

import "github.com/CosmWasm/wasmvm/types"

type SudoMsg struct {
	Initialization  *Initialization  `json:"initialization,omitempty"`
	Authorization   *Authorization   `json:"authorization,omitempty"`
	PreTransaction  *PreTransaction  `json:"pre_transaction,omitempty"`
	PostTransaction *PostTransaction `json:"post_transaction,omitempty"`
}

type Initialization struct {
	Sender  string `json:"sender"`
	Account string `json:"account"`
	Msg     []byte `json:"msg"`
}

type Authorization struct {
	Sender      string `json:"sender"`
	Account     string `json:"account"`
	Signature   []byte `json:"signature"`
	SignedBytes []byte `json:"signed_bytes"`
	Data        []byte `json:"data"`
}

type PreTransaction struct {
	Sender   string            `json:"sender"`
	Account  string            `json:"account"`
	Messages []types.CosmosMsg `json:"msgs"`
}

type PostTransaction struct {
	Sender  string            `json:"sender"`
	Account string            `json:"account"`
	Msgs    []types.CosmosMsg `json:"msgs"`
}
