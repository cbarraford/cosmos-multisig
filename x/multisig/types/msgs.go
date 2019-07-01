package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgCreateWallet defines a CreateWallet message
type MsgCreateWallet struct {
	Name     string           `json:"name"`
	PubKeys  []string         `json:"pub_keys"`
	MinSigTx int              `json:"min_sig_tx"`
	Signers  []sdk.AccAddress `json:"owners"`
	Address  sdk.AccAddress   `json:"address"`
}

// NewMsgCreateWallet is a constructor function for MsgCreateWallet
func NewMsgCreateWallet(name string, pubKeys []string, min int) MsgCreateWallet {
	return MsgCreateWallet{
		Name:     name,
		PubKeys:  pubKeys,
		MinSigTx: min,
	}
}

// Route should return the name of the module
func (msg MsgCreateWallet) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateWallet) Type() string { return "set_wallet" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateWallet) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}
	if len(msg.PubKeys) < msg.MinSigTx {
		return sdk.ErrUnknownRequest("Minimum signature transaction number cannot be larger than the number of public keys")
	}
	if msg.MinSigTx < 1 {
		return sdk.ErrUnknownRequest("Must require at least 1 signatures")
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateWallet) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateWallet) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}
