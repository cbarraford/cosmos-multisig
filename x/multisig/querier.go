package multisig

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the multisig Querier
const (
	ListWallets      = "listWallets"
	GetWallet        = "getWallet"
	ListTransactions = "listTransactions"
	GetTransaction   = "getTransaction"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case ListWallets:
			return queryWallets(ctx, path[1:], req, keeper)
		case GetWallet:
			return getWallet(ctx, path[1:], req, keeper)
		case ListTransactions:
			return queryTransactions(ctx, path[1:], req, keeper)
		case GetTransaction:
			return getTransaction(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown multisig query endpoint")
		}
	}
}

func queryWallets(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var walletList QueryWallets

	iterator := keeper.GetIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		if strings.HasPrefix(string(iterator.Key()), "wallet-") {
			address := strings.TrimPrefix(string(iterator.Key()), "wallet-")
			wallet := keeper.GetWallet(ctx, address)
			for _, pubkey := range wallet.PubKeys {
				if pubkey == path[0] {
					walletList = append(walletList, wallet)
					break
				}
			}

		}
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, walletList)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func getWallet(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	wallet := keeper.GetWallet(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, wallet)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryTransactions(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var transactionList QueryTransactions

	iterator := keeper.GetIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		if strings.HasPrefix(string(iterator.Key()), "transaction-") {
			uidStr := strings.TrimPrefix(string(iterator.Key()), "transaction-")
			transaction := keeper.GetTransaction(ctx, uidStr)
			if transaction.From.String() == path[0] {
				transactionList = append(transactionList, transaction)
			}

		}
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, transactionList)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func getTransaction(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	wallet := keeper.GetTransaction(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, wallet)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
