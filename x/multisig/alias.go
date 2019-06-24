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
	NewMsgBuyName     = types.NewMsgBuyName
	NewMsgSetName     = types.NewMsgSetName
	NewWhois          = types.NewWhois
	NewMultiSigWallet = types.NewMultiSigWallet
	ModuleCdc         = types.ModuleCdc
	RegisterCodec     = types.RegisterCodec
)

type (
	MsgSetName      = types.MsgSetName
	MsgBuyName      = types.MsgBuyName
	QueryResResolve = types.QueryResResolve
	QueryResNames   = types.QueryResNames
	Whois           = types.Whois
	MultiSigWallet  = types.MultiSigWallet
)
