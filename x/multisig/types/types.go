package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// Returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, w.Owner, w.Value, w.Price))
}

// MultiSigWallet is a struct that contains all the metadata of a multiple
// signature wallet
type MultiSigWallet struct {
	Name     string         `json:"name"`
	MinSigTx int            `json:"min_sig_tx"`
	Address  sdk.AccAddress `json:"address"`
	Keys     []string       `json:"pub_keys"`
}

// Returns a new MultiSigWallet with the minprice as the price
func NewMultiSigWallet(name string, keys []string, min int) MultiSigWallet {
	return MultiSigWallet{
		Name:     name,
		MinSigTx: min,
		Keys:     keys,
	}
}

// implement fmt.Stringer
func (w MultiSigWallet) String() string {
	return strings.TrimSpace(
		fmt.Sprintf(
			`Wallet: %s (%d of %d): %s`, w.Name, w.MinSigTx, len(w.Keys), w.Address,
		),
	)
}
