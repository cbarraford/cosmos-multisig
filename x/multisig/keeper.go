package multisig

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/google/uuid"

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

func (k Keeper) GetWallet(ctx sdk.Context, address string) MultiSigWallet {
	address = fmt.Sprintf("wallet-%s", address)
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(address)) {
		return MultiSigWallet{}
	}
	bz := store.Get([]byte(address))
	var wallet MultiSigWallet
	k.cdc.MustUnmarshalBinaryBare(bz, &wallet)
	return wallet
}

// Sets the entire wallet metadata struct for a multisig wallet
func (k Keeper) SetWallet(ctx sdk.Context, wallet MultiSigWallet) {
	address := fmt.Sprintf("wallet-%s", wallet.Address.String())
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(address), k.cdc.MustMarshalBinaryBare(wallet))
}

func (k Keeper) GetTransaction(ctx sdk.Context, uid string) Transaction {
	key := fmt.Sprintf("transaction-%s", uid)
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return Transaction{}
	}
	bz := store.Get([]byte(key))
	var transaction Transaction
	k.cdc.MustUnmarshalBinaryBare(bz, &transaction)
	return transaction
}

func (k Keeper) SetTransaction(ctx sdk.Context, transaction Transaction) {
	key := fmt.Sprintf("transaction-%s", transaction.UUID)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(key), k.cdc.MustMarshalBinaryBare(transaction))
}

func (k Keeper) GetIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

func (k Keeper) DeleteTransaction(ctx sdk.Context, uid uuid.UUID) {
	key := fmt.Sprintf("transaction-%s", uid)
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(key))
}

func (k Keeper) CleanupTransactions(ctx sdk.Context, maxAge int64) {

	currentHeight := ctx.BlockHeight()
	iterator := k.GetIterator(ctx)
	store := ctx.KVStore(k.storeKey)

	for ; iterator.Valid(); iterator.Next() {
		if strings.HasPrefix(string(iterator.Key()), "transaction-") {
			bz := store.Get(iterator.Key())
			var transaction Transaction
			k.cdc.MustUnmarshalBinaryBare(bz, &transaction)
			if currentHeight > (transaction.CreatedAt + maxAge) {
				store.Delete(iterator.Key())
			}
		}
	}
}
