package multisig

import (
	"github.com/cbarraford/parsec"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper parsec.Bank

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(coinKeeper parsec.Bank, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func (k Keeper) GetWallet(ctx parsec.Context, name string) MultiSigWallet {
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
func (k Keeper) SetWallet(ctx parsec.Context, name string, wallet MultiSigWallet) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(wallet))
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
