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
	NewMsgCreateWallet      = types.NewMsgCreateWallet
	NewMultiSigWallet       = types.NewMultiSigWallet
	NewMsgCreateTransaction = types.NewMsgCreateTransaction
	NewTransaction          = types.NewTransaction
	ModuleCdc               = types.ModuleCdc
	RegisterCodec           = types.RegisterCodec
)

type (
	MsgCreateWallet      = types.MsgCreateWallet
	MsgCreateTransaction = types.MsgCreateTransaction
	QueryResResolve      = types.QueryResResolve
	QueryResNames        = types.QueryResNames
	Transaction          = types.Transaction
	MultiSigWallet       = types.MultiSigWallet
)
