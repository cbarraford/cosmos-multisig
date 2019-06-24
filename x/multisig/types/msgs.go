package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgSetName defines a SetName message
type MsgSetWallet struct {
	Name     string         `json:"name"`
	Keys     []string       `json:"keys"`
	MinSigTx int            `json:"min_sig_tx"`
	Address  sdk.AccAddress `json:"address"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgSetWallet(name string, keys []string, min int, address sdk.AccAddress) MsgSetWallet {
	return MsgSetWallet{
		Name:     name,
		Keys:     keys,
		MinSigTx: min,
		Address:  address,
	}
}

// Route should return the name of the module
func (msg MsgSetWallet) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetWallet) Type() string { return "set_wallet" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetWallet) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}
	if len(msg.Keys) < 2 {
		return sdk.ErrUnknownRequest("Must have at least 2 keys")
	}
	if msg.MinSigTx < 2 {
		return sdk.ErrUnknownRequest("Must require at least 2 signatures")
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetWallet) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetWallet) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}
