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
	NewMsgSetWallet   = types.NewMsgSetWallet
	NewMultiSigWallet = types.NewMultiSigWallet
	ModuleCdc         = types.ModuleCdc
	RegisterCodec     = types.RegisterCodec
)

type (
	MsgSetWallet    = types.MsgSetWallet
	QueryResResolve = types.QueryResResolve
	QueryResNames   = types.QueryResNames
	MultiSigWallet  = types.MultiSigWallet
)
