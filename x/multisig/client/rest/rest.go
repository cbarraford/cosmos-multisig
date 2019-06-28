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
	"github.com/cosmos/cosmos-sdk/client/flags"
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
	r.HandleFunc(fmt.Sprintf("/%s/tx", storeName), createUnsignedTransactionHandler(cliCtx)).Methods("PUT")
}

type createUnsignedTransaction struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}

type fee struct {
	Amount []int  `json:"amount"`
	Gas    string `json:"gas"`
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
	Msgs       []msg    `json:"msg"`
	Fee        fee      `json:"fee"`
	Signatures []string `json:"signatures"`
	Memo       string   `json:"memo"`
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

		/*
			baseReq := req.BaseReq.Sanitize()
			if !baseReq.ValidateBasic(w) {
				return
			}
		*/

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
				Fee: fee{
					Amount: make([]int, 0),
					Gas:    fmt.Sprintf("%d", flags.DefaultGasLimit), // hard coded to default gas amount
				},
				Signatures: nil,
				Memo:       "", // TODO: add memo support
			},
		}

		/*
			msg := types.NewMsgSend(from, to, coins)

			stdSignMsg, err := txBldr.BuildSignMsg([]sdk.Msg{msg})
			if err != nil {
				return stdTx, nil
			}
			stdTx := authtypes.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, nil, stdSignMsg.Memo)

			j, err := cliCtx.Codec.MarshalJSON(stdTx)
			if err != nil {
				return err
			}
		*/

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
