package types

import "github.com/CosmWasm/wasmvm/types"

type SudoMsg struct {
	Initialization  *Initialization  `json:"initialization,omitempty"`
	Authorization   *Authorization   `json:"authorization,omitempty"`
	PreTransaction  *PreTransaction  `json:"pre_transaction,omitempty"`
	PostTransaction *PostTransaction `json:"post_transaction,omitempty"`
}

type Authorization struct {
	Senders     []string `json:"senders"`
	Account     string   `json:"account"`
	Signatures  [][]byte `json:"signatures"`
	SignedBytes [][]byte `json:"signed_bytes"`
	Data        []byte   `json:"data"`
}

type PreTransaction struct {
	Sender   string            `json:"sender"`
	Account  string            `json:"account"`
	Messages []types.CosmosMsg `json:"msgs"`
}

type PostTransaction struct {
	Sender   string            `json:"sender"`
	Account  string            `json:"account"`
	Messages []types.CosmosMsg `json:"msgs"`
}
