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
		GetCmdSetWallet(cdc),
	)...)

	return multisigTxCmd
}

// GetCmdSetWallet is the CLI command for sending a SetWallet transaction
func GetCmdSetWallet(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-wallet [name] [min-signatures-required] [pub-keys]",
		Short: "create a new multi-signature wallet",
		Args:  cobra.ExactArgs(2),
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

			msg, err := types.NewMsgSetWallet(args[0], pubKeys, int(minSigs))
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
