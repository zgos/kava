package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/kava-labs/kava/x/auction/types"
)

// GetTxCmd returns the transaction cli commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "transaction commands for the auction module",
	}

	cmds := []*cobra.Command{
		GetCmdPlaceBid(),
		GetCmdPlaceBid2(),
	}

	for _, cmd := range cmds {
		flags.AddTxFlagsToCmd(cmd)
	}

	txCmd.AddCommand(cmds...)

	return txCmd
}

// GetCmdPlaceBid cli command for placing bids on auctions
func GetCmdPlaceBid() *cobra.Command {
	return &cobra.Command{
		Use:     "bid [auction-id] [amount]",
		Short:   "place a bid on an auction",
		Long:    "Place a bid on any type of auction, updating the latest bid amount to [amount]. Collateral auctions must be bid up to their maxbid before entering reverse phase.",
		Example: fmt.Sprintf("  $ %s tx %s bid 34 1000usdx --from myKeyName", version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("auction-id '%s' not a valid uint", args[0])
			}

			amt, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgPlaceBid(id, clientCtx.GetFromAddress().String(), amt)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}

// GetCmdPlaceBid cli command for placing bids on auctions
func GetCmdPlaceBid2() *cobra.Command {
	return &cobra.Command{
		Use:     "bid2 [auction-id] [amount]",
		Short:   "place a bid on an auction",
		Long:    "Place a bid on any type of auction, updating the latest bid amount to [amount]. Collateral auctions must be bid up to their maxbid before entering reverse phase.",
		Example: fmt.Sprintf("  $ %s tx %s bid2 34 1000usdx 1010usdx --from myKeyName", version.AppName, types.ModuleName),
		Args:    cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("auction-id '%s' not a valid uint", args[0])
			}

			var msgs []sdk.Msg
			for _, v := range args[1:] {
				amt, err := sdk.ParseCoinNormalized(v)
				if err != nil {
					return err
				}

				msg := types.NewMsgPlaceBid(id, clientCtx.GetFromAddress().String(), amt)
				err = msg.ValidateBasic()
				if err != nil {
					return err
				}
				msgs = append(msgs, &msg)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)
		},
	}
}
