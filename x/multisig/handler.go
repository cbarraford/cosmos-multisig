package multisig

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetWallet:
			return handleMsgSetWallet(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleMsgSetWallet(ctx sdk.Context, keeper Keeper, msg MsgSetWallet) sdk.Result {
	var err error
	// check the wallet does not already exist
	wallet := keeper.GetWallet(ctx, msg.Name)
	if !wallet.Address.Empty() {
		return sdk.ErrUnauthorized("Wallet already exists").Result()
	}
	wallet, err = NewMultiSigWallet(msg.Name, msg.PubKeys, msg.MinSigTx)
	if err != nil {
		return sdk.ErrUnknownRequest(
			fmt.Sprintf("Error creating new wallet: %s", err.Error()),
		).Result()
	}
	keeper.SetWallet(ctx, msg.Name, wallet)
	return sdk.Result{}
}
