package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"

	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/wallet", storeName), createWalletHandler(cliCtx)).Methods("POST")
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
