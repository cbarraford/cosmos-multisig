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
	)...)

	return multisigTxCmd
}

// GetCmdCreateWallet is the CLI command for sending a CreateWallet transaction
func GetCmdCreateWallet(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-wallet [name] [min-signatures-required] [pub-keys]",
		Short: "create a new multi-signature wallet",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			minSigs, err := strconv.ParseInt(args[1], 0, 64)
			if err != nil {
				return err
			}

			pubKeys := strings.Split(args[2], ",")

			msg := types.NewMsgCreateWallet(args[0], pubKeys, int(minSigs))
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
		Use:   "create-transaction [from] [to] [coins]",
		Short: "create a new multi-signature transaction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

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

			msg := types.NewMsgCreateTransaction(from, to, coins)
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
