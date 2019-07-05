package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	mtypes "github.com/cbarraford/cosmos-multisig/x/multisig/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/gorilla/mux"
)

const (
	walletAddress = "address"
	transactionID = "transaction_id"
)

var (
	ModuleCdc = mtypes.ModuleCdc
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/wallet/{%s}", storeName, walletAddress), getWalletHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/transaction/{%s}", storeName, transactionID), getTransactionHandler(cliCtx, storeName)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/wallets/{%s}", storeName, walletAddress), walletsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/transactions/{%s}", storeName, transactionID), transactionsHandler(cliCtx, storeName)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/wallet", storeName), createWalletHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/transaction", storeName), createTransactionHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/transaction/sign", storeName), signTransactionHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/transaction/complete", storeName), completeTransactionHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/tx", storeName), createUnsignedTransactionHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/sign/multi", storeName), multiSignHandler(cliCtx)).Methods("POST")

	r.HandleFunc(fmt.Sprintf("/%s/broadcast", storeName), broadcastTxRequest(cliCtx)).Methods("POST")
}

func getWalletHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[walletAddress]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/getWallet/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getTransactionHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[transactionID]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/getTransaction/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func walletsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[walletAddress]
		log.Printf("Param %+v", paramType)

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/listWallets/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func transactionsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[transactionID]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/listTransactions/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
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

type createTransaction struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Amount  sdk.Int        `json:"amount"`
	Denom   string         `json:"denom"`
	Signers []string       `json:"signers"`
}

func createTransactionHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createTransaction
		var err error

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		signers := make([]sdk.AccAddress, len(req.Signers))
		for i, _ := range req.Signers {
			signers[i], err = sdk.AccAddressFromBech32(req.Signers[i])
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		// create the message
		msg := mtypes.NewMsgCreateTransaction(req.From, req.To, req.Amount, req.Denom, signers)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type signTransaction struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	UUID      string       `json:"uuid"`
	Signature string       `json:"signature"`
	PubKey    string       `json:"pub_key"`
}

func signTransactionHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signTransaction
		var err error

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			// TODO: is this needed?
			// return
		}

		pubkey, err := sdk.GetAccPubKeyBech32(req.PubKey)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		sig := mtypes.Signature{
			PubKey:    pubkey,
			Signature: req.Signature,
		}

		// create the message
		msg := mtypes.NewMsgSignTransaction(req.UUID, sig)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type completeTransaction struct {
	BaseReq rest.BaseReq `json:"base_req"`
	UUID    string       `json:"uuid"`
	TxID    string       `json:"tx_id"`
}

func completeTransactionHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req completeTransaction
		var err error

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			// TODO: is this needed?
			// return
		}

		// create the message
		msg := mtypes.NewMsgCompleteTransaction(req.UUID, req.TxID)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type createWallet struct {
	Name     string       `json:"name"`
	BaseReq  rest.BaseReq `json:"base_req"`
	Address  string       `json:"address"`
	MinSigTx int          `json:"min_sig_tx"`
	PubKeys  []string     `json:"pub_keys"`
	Signers  []string     `json:"signers"`
}

func createWalletHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createWallet
		var err error

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		signers := make([]sdk.AccAddress, len(req.Signers))
		for i, _ := range req.Signers {
			signers[i], err = sdk.AccAddressFromBech32(req.Signers[i])
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		// create the message
		msg := mtypes.NewMsgCreateWallet(req.Name, req.PubKeys, req.MinSigTx, signers)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
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

// BroadcastReq defines a tx broadcasting request.
type BroadcastReq struct {
	Tx   types.StdTx `json:"tx"`
	Mode string      `json:"mode"`
}

// broadcastTxRequest implements a tx broadcasting handler that is responsible
// for broadcasting a valid and signed tx to a full node. The tx can be
// broadcasted via a sync|async|block mechanism.
func broadcastTxRequest(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BroadcastReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			log.Printf("Foo1")
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			log.Printf("Foo2")
			return
		}

		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(req.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			log.Printf("Foo3")
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(req.Mode)

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			log.Printf("Foo4")
			return
		}

		log.Printf("Foo5")
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
