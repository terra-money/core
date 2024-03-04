package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ sdk.Msg = &MsgCreateSmartAccount{}
	_ sdk.Msg = &MsgDisableSmartAccount{}
	_ sdk.Msg = &MsgUpdateAuthorization{}
	_ sdk.Msg = &MsgUpdateTransactionHooks{}
)

func NewMsgCreateSmartAccount(account string) *MsgCreateSmartAccount {
	return &MsgCreateSmartAccount{
		Account: account,
	}
}

func NewMsgDisableSmartAccount(account string) *MsgDisableSmartAccount {
	return &MsgDisableSmartAccount{
		Account: account,
	}
}

func NewMsgUpdateAuthorization(account string, authorizationMsgs []*AuthorizationMsg, fallback bool) *MsgUpdateAuthorization {
	return &MsgUpdateAuthorization{
		Account:           account,
		AuthorizationMsgs: authorizationMsgs,
		Fallback:          fallback,
	}
}

// GetSignBytes implements the LegacyMsg interface.
func (m MsgCreateSmartAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m MsgCreateSmartAccount) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Account)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m MsgCreateSmartAccount) ValidateBasic() error {
	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (m MsgDisableSmartAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m MsgDisableSmartAccount) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Account)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m MsgDisableSmartAccount) ValidateBasic() error {
	return nil
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
	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (m MsgUpdateTransactionHooks) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m MsgUpdateTransactionHooks) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Account)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m MsgUpdateTransactionHooks) ValidateBasic() error {
	return nil
}
