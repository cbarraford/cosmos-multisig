package types

import (
	"encoding/hex"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MultiSigWallet is a struct that contains all the metadata of a multiple
// signature wallet
type MultiSigWallet struct {
	Name     string         `json:"name"`
	MinSigTx int            `json:"min_sig_tx"`
	Address  sdk.AccAddress `json:"address"`
	PubKeys  []string       `json:"pub_keys"`
}

// Returns a new MultiSigWallet with the minprice as the price
func NewMultiSigWallet(name string, pubKeys []string, min int) (MultiSigWallet, error) {
	// encode name into a hex []byte
	src := []byte(name)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	addr, err := sdk.AccAddressFromHex(string(dst))
	if err != nil {
		return MultiSigWallet{}, err
	}

	// check if we need to set a default min keys
	if min < 2 {
		min = len(pubKeys) - 1
	}

	return MultiSigWallet{
		Name:     name,
		MinSigTx: min,
		PubKeys:  pubKeys,
		Address:  addr,
	}, nil
}

// implement fmt.Stringer
func (w MultiSigWallet) String() string {
	return strings.TrimSpace(
		fmt.Sprintf(
			`Wallet: %s (%d of %d): %s`, w.Name, w.MinSigTx, len(w.PubKeys), w.Address,
		),
	)
}
