package multisig

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "multisig" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateWallet:
			return handleMsgCreateWallet(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized multisig Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to create wallet
func handleMsgCreateWallet(ctx sdk.Context, keeper Keeper, msg MsgCreateWallet) sdk.Result {
	var err error
	// check the wallet does not already exist
	wallet := keeper.GetWallet(ctx, msg.Address.String())
	if !wallet.Address.Empty() {
		return sdk.ErrUnauthorized("Wallet already exists").Result()
	}
	wallet, err = NewMultiSigWallet(msg.Name, msg.PubKeys, msg.MinSigTx)
	if err != nil {
		return sdk.ErrUnknownRequest(
			fmt.Sprintf("Error creating new wallet: %s", err.Error()),
		).Result()
	}
	keeper.CreateWallet(ctx, wallet)
	return sdk.Result{}
}

// Handle a message to create transaction
func handleMsgCreateTransaction(ctx sdk.Context, keeper Keeper, msg MsgCreateTransaction) sdk.Result {
	// check the transaction does not already exist
	transaction := keeper.GetTransaction(ctx, msg.UUID)
	if !transaction.From.Empty() {
		return sdk.ErrUnauthorized("Transaction already exists").Result()
	}
	transaction = NewTransaction(msg.From, msg.To, msg.Coins, ctx.BlockHeight())
	keeper.CreateTransaction(ctx, transaction)
	return sdk.Result{}
}
