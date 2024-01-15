package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgAuctionBid{}
)

// GetSignBytes implements the LegacyMsg interface.
func (m MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m MsgUpdateParams) ValidateBasic() error {
	return nil
}

func NewMsgAuctionBid(bidder sdk.AccAddress, bid sdk.Coin, transactions [][]byte) *MsgAuctionBid {
	return &MsgAuctionBid{
		Bidder:       bidder.String(),
		Bid:          bid,
		Transactions: transactions,
	}
}

// GetSignBytes implements the LegacyMsg interface.
func (m MsgAuctionBid) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgAuctionBid message.
func (m MsgAuctionBid) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Bidder)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m MsgAuctionBid) ValidateBasic() error {
	return nil
}
