package cli

import (
	"github.com/spf13/cobra"

	"github.com/cbarraford/cosmos-multisig/x/multisig/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	multisigTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Multisig transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	multisigTxCmd.AddCommand(client.PostCommands()...)

	return multisigTxCmd
}
