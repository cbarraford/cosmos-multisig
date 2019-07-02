package multisig

import (
	"github.com/cbarraford/cosmos-multisig/x/multisig/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewMsgCreateWallet        = types.NewMsgCreateWallet
	NewMultiSigWallet         = types.NewMultiSigWallet
	NewMsgCreateTransaction   = types.NewMsgCreateTransaction
	NewMsgSignTransaction     = types.NewMsgSignTransaction
	NewMsgCompleteTransaction = types.NewMsgCompleteTransaction
	NewTransaction            = types.NewTransaction
	ModuleCdc                 = types.ModuleCdc
	RegisterCodec             = types.RegisterCodec
)

type (
	MsgCreateWallet        = types.MsgCreateWallet
	MsgCreateTransaction   = types.MsgCreateTransaction
	MsgSignTransaction     = types.MsgSignTransaction
	MsgCompleteTransaction = types.MsgCompleteTransaction
	QueryWallets           = types.QueryWallets
	QueryTransactions      = types.QueryTransactions
	Transaction            = types.Transaction
	Signature              = types.Signature
	MultiSigWallet         = types.MultiSigWallet
)
