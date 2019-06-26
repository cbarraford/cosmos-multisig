package multisig

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the multisig Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func (k Keeper) GetWallet(ctx sdk.Context, name string) MultiSigWallet {
	name = fmt.Sprintf("wallet-%s", name)
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(name)) {
		return MultiSigWallet{}
	}
	bz := store.Get([]byte(name))
	var wallet MultiSigWallet
	k.cdc.MustUnmarshalBinaryBare(bz, &wallet)
	return wallet
}

// Sets the entire wallet metadata struct for a multisig wallet
func (k Keeper) SetWallet(ctx sdk.Context, name string, wallet MultiSigWallet) {
	name = fmt.Sprintf("wallet-%s", name)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(wallet))
}

func (k Keeper) GetWalletIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
