package cli

import (
	"strconv"
	"strings"

	"github.com/cbarraford/cosmos-multisig/x/multisig/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	multisigTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Multisig transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	multisigTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateWallet(cdc),
		GetCmdCreateTransaction(cdc),
		GetCmdSignTransaction(cdc),
		GetCmdCompleteTransaction(cdc),
	)...)

	return multisigTxCmd
}

// GetCmdCreateWallet is the CLI command for sending a CreateWallet transaction
func GetCmdCreateWallet(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-wallet [name] [min-signatures-required] [pub-keys], [addresses]",
		Short: "create a new multi-signature wallet",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			minSigs, err := strconv.ParseInt(args[1], 0, 64)
			if err != nil {
				return err
			}

			pubKeys := strings.Split(args[2], ",")

			addrs := strings.Split(args[3], ",")
			signers := make([]sdk.AccAddress, len(addrs))
			for i, addr := range addrs {
				signers[i], err = sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgCreateWallet(args[0], pubKeys, int(minSigs), signers)
			if err != nil {
				return err
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdCreateTransaction is the CLI command for sending a CreateTransaction transaction
func GetCmdCreateTransaction(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-transaction [from] [to] [coins] [signers]",
		Short: "create a new multi-signature transaction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			from, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			addrs := strings.Split(args[3], ",")
			signers := make([]sdk.AccAddress, len(addrs))
			for i, addr := range addrs {
				signers[i], err = sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgCreateTransaction(from, to, coins[0].Amount, coins[0].Denom, signers)
			if err != nil {
				return err
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdSignTransaction is the CLI command for saving a transaction signature
func GetCmdSignTransaction(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "save-transaction-signature [uuid] [pubkey] [signature] [signers]",
		Short: "Save a signature generated for a specific transaction",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			addrs := strings.Split(args[3], ",")
			signers := make([]sdk.AccAddress, len(addrs))
			for i, addr := range addrs {
				signers[i], err = sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgSignTransaction(args[0], args[1], args[2], signers)
			if err != nil {
				return err
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdCompleteTransaction is the CLI command for saving a transaction signature
func GetCmdCompleteTransaction(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "complete-transaction [uuid] [transaction_id] [signers]",
		Short: "Save a blockchain transaction id to a transaction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			addrs := strings.Split(args[3], ",")
			signers := make([]sdk.AccAddress, len(addrs))
			for i, addr := range addrs {
				signers[i], err = sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgCompleteTransaction(args[0], args[1], signers)
			if err != nil {
				return err
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
