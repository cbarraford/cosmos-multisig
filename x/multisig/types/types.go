package types

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
)

// MultiSigWallet is a struct that contains all the metadata of a multiple
// signature wallet
type MultiSigWallet struct {
	Name     string         `json:"name"`
	MinSigTx int            `json:"min_sig_tx"`
	Address  sdk.AccAddress `json:"address"`
	PubKeys  []string       `json:"pub_keys"`
}

func createAddress(name string) (sdk.AccAddress, error) {
	// encode name into a hex []byte
	src := []byte(name)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	return sdk.AccAddressFromHex(string(dst))
}

func validateMultisigThreshold(k, nKeys int) error {
	if k <= 0 {
		return fmt.Errorf("threshold must be a positive integer")
	}
	if nKeys < k {
		return fmt.Errorf(
			"threshold k of n multisignature: %d < %d", nKeys, k)
	}
	return nil
}

// Returns a new MultiSigWallet with the minprice as the price
func NewMultiSigWallet(name string, pubKeys []string, min int) (MultiSigWallet, error) {
	var err error

	err = validateMultisigThreshold(min, len(pubKeys))
	if err != nil {
		return MultiSigWallet{}, err
	}

	cryptoPubKeys := make([]crypto.PubKey, len(pubKeys))
	for i, _ := range cryptoPubKeys {
		cryptoPubKeys[i], err = sdk.GetAccPubKeyBech32(pubKeys[i])
		if err != nil {
			return MultiSigWallet{}, err
		}
	}

	multikey := multisig.NewPubKeyMultisigThreshold(min, cryptoPubKeys)
	info := keys.NewMultiInfo("multisig", multikey)

	return MultiSigWallet{
		Name:     name,
		MinSigTx: min,
		PubKeys:  pubKeys,
		Address:  info.GetAddress(),
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

type Signature struct {
	PubKey    crypto.PubKey `json:"pub_key"`
	Signature string        `json:"signature"`
}

type Transaction struct {
	UUID       uuid.UUID      `json:"uuid"`
	From       sdk.AccAddress `json:"from_address"`
	To         sdk.AccAddress `json:"to_address"`
	Coins      sdk.Coins      `json:"coins"`
	Signatures []Signature    `json:"signatures"`
	CreatedAt  int64          `json:"created_at"` // block height
}

func NewTransaction(from, to sdk.AccAddress, coins sdk.Coins, height int64) Transaction {
	return Transaction{
		UUID:      uuid.New(),
		From:      from,
		To:        to,
		Coins:     coins,
		CreatedAt: height,
	}
}

func (t Transaction) String() string {
	return strings.TrimSpace(
		fmt.Sprintf(
			`Transaction (%s): %s --> %s %+v`, t.UUID, t.From, t.To, t.Coins,
		),
	)
}
