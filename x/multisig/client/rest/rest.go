package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	mtypes "github.com/cbarraford/cosmos-multisig/x/multisig/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"

	"github.com/gorilla/mux"
)

var (
	ModuleCdc = mtypes.ModuleCdc
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/wallet", storeName), createWalletHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/tx", storeName), createUnsignedTransactionHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/sign/multi", storeName), multiSignHandler(cliCtx)).Methods("POST")
}

type pubkey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type signature struct {
	PubKey    pubkey `json:"pub_key,omitempty"`
	Signature string `json:"signature"`
}

type multiSign struct {
	PubKeys       []string            `json:"pub_keys"`
	MinSigTx      int                 `json:"min_sig_tx"`
	Signatures    []signature         `json:"signatures"`
	Tx            unsignedTransaction `json:"unsigned_tx"`
	ChainID       string              `json:"chain_id"`
	AccountNumber uint64              `json:"account_number"`
	Sequence      uint64              `json:"sequence"`
}

func multiSignHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req multiSign
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		// Important: The input signature must be in order relative to their
		// public key was in the list of public keys were used to create the
		// multisig wallet.
		signatures := make([]string, len(req.Signatures))
		for i, stdSig := range req.Signatures {
			if i < (len(req.Signatures) - 1) {
				sig := []byte(stdSig.Signature)
				// since we are not the last signature, make edits...
				// remove the '==' at the end of the string
				sig = sig[:len(sig)-2]
				// increment the last character by 1
				sig[len(sig)-1] += 1
				signatures[i] = string(sig)
			} else {
				signatures[i] = stdSig.Signature
			}
		}

		// prepend base string "CgUIAxIB4B"
		signatures = append([]string{"CgUIAxIB4B"}, signatures...)
		multisignature := strings.Join(signatures[:], "JA")

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, fmt.Sprintf(`{"signature":"%s"}`, multisignature))
	}
}

type createUnsignedTransaction struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	Memo   string `json:"memo"`
}

type msg struct {
	Type  string `json:"type"`
	Value value2 `json:"value"`
}

type value2 struct {
	FromAddress string    `json:"from_address"`
	ToAddress   string    `json:"to_address"`
	Amount      sdk.Coins `json:"amount"`
}

type value1 struct {
	Msgs       []msg        `json:"msg"`
	Fee        types.StdFee `json:"fee"`
	Signatures []string     `json:"signatures"`
	Memo       string       `json:"memo"`
}

type unsignedTransaction struct {
	Type  string `json:"type"`
	Value value1 `json:"value"`
}

func createUnsignedTransactionHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createUnsignedTransaction
		var err error

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		_, err = sdk.AccAddressFromBech32(req.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		_, err = sdk.AccAddressFromBech32(req.To)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		coins, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx := unsignedTransaction{
			Type: "auth/StdTx",
			Value: value1{
				Msgs: []msg{
					{
						Type: "cosmos-sdk/MsgSend",
						Value: value2{
							FromAddress: req.From,
							ToAddress:   req.To,
							Amount:      coins,
						},
					},
				},
				Fee:        types.NewStdFee(flags.DefaultGasLimit, sdk.Coins{}),
				Signatures: nil,
				Memo:       req.Memo,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		j, _ := json.Marshal(tx)
		io.WriteString(w, string(j))
	}
}

type createWallet struct {
	Address  string   `json:"address"`
	MinSigTx int      `json:"min_sig_tx"`
	PubKeys  []string `json:"pub_keys"`
}

func createWalletHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createWallet
		var info keys.Info

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Unmarshal
		err = json.Unmarshal(b, &req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = validateMultisigThreshold(req.MinSigTx, len(req.PubKeys))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		sort.Slice(req.PubKeys, func(i, j int) bool {
			return bytes.Compare([]byte(req.PubKeys[i]), []byte(req.PubKeys[j])) < 0
		})

		pubKeys := make([]crypto.PubKey, len(req.PubKeys))
		for i, _ := range req.PubKeys {
			var err error
			pubKeys[i], err = sdk.GetAccPubKeyBech32(req.PubKeys[i])
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		multikey := multisig.NewPubKeyMultisigThreshold(req.MinSigTx, pubKeys)
		info = keys.NewMultiInfo("multi", multikey)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		req.Address = fmt.Sprintf("%s", info.GetAddress())
		wallet, _ := json.Marshal(req)
		io.WriteString(w, string(wallet))
	}
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
