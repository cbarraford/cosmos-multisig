package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgCreateWallet defines a CreateWallet message
type MsgCreateWallet struct {
	MinSigTx int              `json:"min_sig_tx"`
	Name     string           `json:"name"`
	PubKeys  []string         `json:"pub_keys"`
	Signers  []sdk.AccAddress `json:"signers"`
}

// NewMsgCreateWallet is a constructor function for MsgCreateWallet
func NewMsgCreateWallet(name string, pubKeys []string, min int, signers []sdk.AccAddress) MsgCreateWallet {
	return MsgCreateWallet{
		Name:     name,
		PubKeys:  pubKeys,
		MinSigTx: min,
		Signers:  signers,
	}
}

// Route should return the name of the module
func (msg MsgCreateWallet) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateWallet) Type() string { return "set_wallet" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateWallet) ValidateBasic() sdk.Error {
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
	return msg.Signers
}

// MsgCreateTransaction defines a CreateTransaction message
type MsgCreateTransaction struct {
	Amount  sdk.Int          `json:"amount"`
	Denom   string           `json:"denom"`
	From    sdk.AccAddress   `json:"from_address"`
	Signers []sdk.AccAddress `json:"signers"`
	To      sdk.AccAddress   `json:"to_address"`
	UUID    string           `json:"uuid"`
}

// NewMsgCreateTransaction is a constructor function for MsgCreateTransaction
func NewMsgCreateTransaction(from, to sdk.AccAddress, amount sdk.Int, denom string, signers []sdk.AccAddress) MsgCreateTransaction {
	return MsgCreateTransaction{
		UUID:    uuid.New().String(),
		From:    from,
		To:      to,
		Amount:  amount,
		Denom:   denom,
		Signers: signers,
	}
}

// Route should return the name of the module
func (msg MsgCreateTransaction) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreateTransaction) Type() string { return "create_transaction" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateTransaction) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress(msg.From.String())
	}
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress(msg.To.String())
	}
	/*
		if msg.Coins.IsZero() {
			return sdk.ErrUnknownRequest("Cannot have zero coins")
		}
		if msg.Coins.IsValid() {
			// TODO: readd this back at some point
			//return sdk.ErrUnknownRequest("Coins must be valid")
		}
	*/
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreateTransaction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreateTransaction) GetSigners() []sdk.AccAddress {
	return msg.Signers
}

// MsgSignTransaction defines a SignTransaction message
type MsgSignTransaction struct {
	PubKey       string           `json:"pubkey"`
	PubKeyBase64 string           `json:"pubkey_base64"`
	Signature    string           `json:"signature"`
	Signers      []sdk.AccAddress `json:"signers"`
	UUID         string           `json:"uuid"`
}

// NewMsgSignTransaction is a constructor function for MsgCreateTransaction
func NewMsgSignTransaction(uid, pubkey, pubkey_base64, sig string, signers []sdk.AccAddress) MsgSignTransaction {
	return MsgSignTransaction{
		UUID:         uid,
		PubKey:       pubkey,
		PubKeyBase64: pubkey_base64,
		Signature:    sig,
		Signers:      signers,
	}
}

// Route should return the name of the module
func (msg MsgSignTransaction) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSignTransaction) Type() string { return "sign_transaction" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSignTransaction) ValidateBasic() sdk.Error {
	if len(msg.UUID) == 0 {
		return sdk.ErrUnknownRequest("UUID cannot be blank")
	}
	if msg.PubKey == "" {
		return sdk.ErrUnknownRequest("Pubkey cannot be blank")
	}
	if msg.PubKeyBase64 == "" {
		return sdk.ErrUnknownRequest("Pubkey (base64) cannot be blank")
	}
	if msg.Signature == "" {
		return sdk.ErrUnknownRequest("Signature cannot be blank")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSignTransaction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSignTransaction) GetSigners() []sdk.AccAddress {
	return msg.Signers
}

// MsgCompleteTransaction defines complete a transaction
type MsgCompleteTransaction struct {
	Signers []sdk.AccAddress `json:"signers"`
	TxID    string           `json:"tx_id"`
	UUID    string           `json:"uuid"`
}

// NewMsgCompleteTransaction is a constructor function for MsgCompleteTransaction
func NewMsgCompleteTransaction(uid, txID string, signers []sdk.AccAddress) MsgCompleteTransaction {
	return MsgCompleteTransaction{
		UUID:    uid,
		TxID:    txID,
		Signers: signers,
	}
}

// Route should return the name of the module
func (msg MsgCompleteTransaction) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCompleteTransaction) Type() string { return "sign_transaction" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCompleteTransaction) ValidateBasic() sdk.Error {
	if len(msg.UUID) == 0 {
		return sdk.ErrUnknownRequest("UUID cannot be blank")
	}
	if msg.TxID == "" {
		return sdk.ErrUnknownRequest("Transaction ID cannot be blank")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCompleteTransaction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCompleteTransaction) GetSigners() []sdk.AccAddress {
	return msg.Signers
}
