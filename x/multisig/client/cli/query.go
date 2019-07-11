package cli

import (
	"fmt"

	"github.com/cbarraford/cosmos-multisig/x/multisig/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	msigQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the multisig module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	msigQueryCmd.AddCommand(client.GetCommands(
		GetCmdWallet(storeKey, cdc),
		GetCmdWallets(storeKey, cdc),
		GetCmdTransaction(storeKey, cdc),
		GetCmdTransactions(storeKey, cdc),
	)...)
	return msigQueryCmd
}

// GetCmdWallet queries information about a domain
func GetCmdWallet(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-wallet [address]",
		Short: "Get wallet by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/getWallet/%s", queryRoute, addr), nil)
			if err != nil {
				fmt.Printf("could not resolve wallet - %s \n", addr)
				return nil
			}

			var out types.MultiSigWallet
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdWallets queries a list of wallets contains a specific public key
func GetCmdWallets(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "query-wallets [pub_key]",
		Short: "Query for a list of wallets by public key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			pubKey := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/listWallets/%s", queryRoute, pubKey), nil)
			if err != nil {
				fmt.Printf("could not get query wallets\n")
				return nil
			}

			var out types.QueryWallets
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdTransaction queries information about a domain
func GetCmdTransaction(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-transaction [uuid]",
		Short: "Get transaction by uuid",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			uid := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/getTransaction/%s", queryRoute, uid), nil)
			if err != nil {
				fmt.Printf("could not resolve transaction - %s \n", uid)
				return nil
			}

			var out types.Transaction
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdTransactions queries a list of transaction for a specific wallet
func GetCmdTransactions(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "query-transactions [wallet_address]",
		Short: "Query for a list of transaction by wallet address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/listTransactions/%s", queryRoute, addr), nil)
			if err != nil {
				fmt.Printf("could not get query wallets\n")
				return nil
			}

			var out types.QueryTransactions
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
