package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateWallet{}, "multisig/CreateWallet", nil)
	cdc.RegisterConcrete(MsgCreateTransaction{}, "multisig/CreateTransaction", nil)
	cdc.RegisterConcrete(MsgSignTransaction{}, "multisig/SignTransaction", nil)
	cdc.RegisterConcrete(MsgCompleteTransaction{}, "multisig/CompleteTransaction", nil)
}
