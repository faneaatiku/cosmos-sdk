package cli

import (
	"strings"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/circuit/types"
	"github.com/spf13/cobra"
)

// NewTxCmd returns a root CLI command handler for all x/circuit transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Circuit transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		AuthorizeCircuitBreakerCmd(),
		TripCircuitBreakerCmd(),
	)

	return txCmd
}

// AuthorizeCircuitBreakerCmd returns a CLI command handler for creating a MsgAuthorizeCircuitBreaker transaction.
func AuthorizeCircuitBreakerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "authorize [granter] [grantee] [permission_level] [limit_type_urls]",
		Short: "Authorize an account to trip the circuit breaker.",
		Long: `Authorize an account to trip the circuit breaker.
		"SOME_MSGS" =     1,
		"ALL_MSGS" =      2,
		"SUPER_ADMIN" =   3,

		Example: 

		<app> circuit authorize [address] [address] 0 "cosmos.bank.v1beta1.MsgSend,cosmos.bank.v1beta1.MsgMultiSend"
		`,
		Args: cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			grantee, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			lvl, err := math.ParseUint(args[2])
			if err != nil {
				return err
			}

			var typeUrls []string
			if len(args) == 4 {
				typeUrls = strings.Split(args[3], ",")
			}

			permission := types.Permissions{Level: types.Permissions_Level(lvl.Uint64()), LimitTypeUrls: typeUrls}

			msg := types.NewMsgAuthorizeCircuitBreaker(string(clientCtx.GetFromAddress()), string(grantee), &permission)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// TripCircuitBreakerCmd returns a CLI command handler for creating a MsgTripCircuitBreaker transaction.
func TripCircuitBreakerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Disable [type_url]",
		Short: "Disable a message from being executed ",
		Long: `Disable a message  from entering the mempool and/or being executed
		
		Example: 

		<app> circuit authorize trip "cosmos.bank.v1beta1.MsgSend,cosmos.bank.v1beta1.MsgMultiSend"
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var msgTypeUrls []string
			if len(args) == 4 {
				msgTypeUrls = strings.Split(args[0], ",")
			}

			msg := types.NewMsgTripCircuitBreaker(string(clientCtx.GetFromAddress()), msgTypeUrls)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// ResetCircuitBreakerCmd returns a CLI command handler for creating a MsgRestCircuitBreaker transaction.
func ResetCircuitBreakerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset [type_url]",
		Short: "Enable a message to be executed ",
		Long: `Enable a message  that was disabled from entering the mempool and/or being executed
		
		Example: 

		<app> circuit authorize reset "cosmos.bank.v1beta1.MsgSend,cosmos.bank.v1beta1.MsgMultiSend"
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var msgTypeUrls []string
			if len(args) == 4 {
				msgTypeUrls = strings.Split(args[0], ",")
			}

			msg := types.NewMsgResetCircuitBreaker(string(clientCtx.GetFromAddress()), msgTypeUrls)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}