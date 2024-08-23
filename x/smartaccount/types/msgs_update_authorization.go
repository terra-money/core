package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgUpdateTransactionHooks{}
)

func NewMsgUpdateAuthorization(account string, authorizationMsgs []*AuthorizationMsg, fallback bool) *MsgUpdateAuthorization {
	return &MsgUpdateAuthorization{
		Account:           account,
		AuthorizationMsgs: authorizationMsgs,
		Fallback:          fallback,
	}
}

// GetSignBytes implements the LegacyMsg interface.
func (m MsgUpdateAuthorization) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m MsgUpdateAuthorization) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Account)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m MsgUpdateAuthorization) ValidateBasic() error {
	if m.Account == "" {
		return sdkerrors.ErrInvalidAddress.Wrap("account cannot be empty")
	}
	for _, auth := range m.AuthorizationMsgs {
		if err := auth.ValidateBasic(); err != nil {
			return err
		}
	}
	return nil
}

func (a AuthorizationMsg) ValidateBasic() error {
	if a.ContractAddress == "" {
		return sdkerrors.ErrInvalidAddress.Wrap("auth contract address cannot be empty")
	}
	if a.InitMsg == nil {
		return sdkerrors.ErrInvalidRequest.Wrap("init msg cannot be nil")
	}
	if a.InitMsg.Account == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("init msg account cannot be empty")
	}
	return nil
}
