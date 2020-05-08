package cli

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/irismod/service/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	serviceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Service transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	serviceTxCmd.AddCommand(flags.PostCommands(
		GetCmdDefineService(cdc),
		GetCmdBindService(cdc),
		GetCmdUpdateServiceBinding(cdc),
		GetCmdSetWithdrawAddr(cdc),
		GetCmdDisableServiceBinding(cdc),
		GetCmdEnableServiceBinding(cdc),
		GetCmdRefundServiceDeposit(cdc),
		GetCmdCallService(cdc),
		GetCmdRespondService(cdc),
		GetCmdPauseRequestContext(cdc),
		GetCmdStartRequestContext(cdc),
		GetCmdKillRequestContext(cdc),
		GetCmdUpdateRequestContext(cdc),
		GetCmdWithdrawEarnedFees(cdc),
	)...)

	return serviceTxCmd
}

// GetCmdDefineService implements defining a service command
func GetCmdDefineService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "define",
		Short: "Define a new service",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Define a new service based on the given params.

Example:
$ %s tx service define --name=<service name> --description=<service description> --author-description=<author description> 
--tags=<tag1,tag2,...> --schemas=<schemas content or path/to/schemas.json> --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			author := cliCtx.GetFromAddress()

			name := viper.GetString(FlagName)
			description := viper.GetString(FlagDescription)
			authorDescription := viper.GetString(FlagAuthorDescription)
			tags := viper.GetStringSlice(FlagTags)
			schemas := viper.GetString(FlagSchemas)

			if !json.Valid([]byte(schemas)) {
				schemasContent, err := ioutil.ReadFile(schemas)
				if err != nil {
					return fmt.Errorf("invalid schemas: neither JSON input nor path to .json file were provided")
				}

				if !json.Valid(schemasContent) {
					return fmt.Errorf("invalid schemas: .json file content is invalid JSON")
				}

				schemas = string(schemasContent)
			}

			buf := bytes.NewBuffer([]byte{})
			if err := json.Compact(buf, []byte(schemas)); err != nil {
				return fmt.Errorf("failed to compact the schemas")
			}

			schemas = buf.String()
			fmt.Printf("schemas content: \n%s\n", schemas)

			msg := types.NewMsgDefineService(name, description, tags, author, authorDescription, schemas)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsDefineService)
	_ = cmd.MarkFlagRequired(FlagName)
	_ = cmd.MarkFlagRequired(FlagSchemas)

	return cmd
}

// GetCmdBindService implements binding a service command
func GetCmdBindService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind",
		Short: "Bind an existing service definition",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Bind an existing service definition.

Example:
$ %s tx service bind --service-name=<service-name> --deposit=1stake 
--pricing=<pricing content or path/to/pricing.json> --qos=50 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			var err error
			var provider sdk.AccAddress

			providerStr := viper.GetString(FlagProvider)
			if len(providerStr) > 0 {
				provider, err = sdk.AccAddressFromBech32(providerStr)
				if err != nil {
					return err
				}
			} else {
				provider = owner
			}

			serviceName := viper.GetString(FlagServiceName)
			qos := uint64(viper.GetInt64(FlagQoS))

			depositStr := viper.GetString(FlagDeposit)
			deposit, err := sdk.ParseCoins(depositStr)
			if err != nil {
				return err
			}

			pricing := viper.GetString(FlagPricing)

			if !json.Valid([]byte(pricing)) {
				pricingContent, err := ioutil.ReadFile(pricing)
				if err != nil {
					return fmt.Errorf("invalid pricing: neither JSON input nor path to .json file were provided")
				}

				if !json.Valid(pricingContent) {
					return fmt.Errorf("invalid pricing: .json file content is invalid JSON")
				}

				pricing = string(pricingContent)
			}

			buf := bytes.NewBuffer([]byte{})
			if err := json.Compact(buf, []byte(pricing)); err != nil {
				return fmt.Errorf("failed to compact the pricing")
			}

			pricing = buf.String()

			msg := types.NewMsgBindService(serviceName, provider, deposit, pricing, qos, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsBindService)
	_ = cmd.MarkFlagRequired(FlagServiceName)
	_ = cmd.MarkFlagRequired(FlagDeposit)
	_ = cmd.MarkFlagRequired(FlagPricing)
	_ = cmd.MarkFlagRequired(FlagQoS)

	return cmd
}

// GetCmdUpdateServiceBinding implements updating a service binding command
func GetCmdUpdateServiceBinding(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-binding [service-name] [provider-address]",
		Short: "Update an existing service binding",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update an existing service binding.

Example:
$ %s tx service update-binding <service-name> <provider-address> --deposit=1stake 
--pricing=<pricing content or path/to/pricing.json> --qos=50 --from mykey
`,
				version.ClientName,
			),
		),
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			var err error
			var provider sdk.AccAddress

			if len(args) > 1 {
				provider, err = sdk.AccAddressFromBech32(args[1])
				if err != nil {
					return err
				}
			} else {
				provider = owner
			}

			var deposit sdk.Coins

			depositStr := viper.GetString(FlagDeposit)
			if len(depositStr) != 0 {
				deposit, err = sdk.ParseCoins(depositStr)
				if err != nil {
					return err
				}
			}

			pricing := viper.GetString(FlagPricing)

			if len(pricing) != 0 {
				if !json.Valid([]byte(pricing)) {
					pricingContent, err := ioutil.ReadFile(pricing)
					if err != nil {
						return fmt.Errorf("invalid pricing: neither JSON input nor path to .json file were provided")
					}

					if !json.Valid(pricingContent) {
						return fmt.Errorf("invalid pricing: .json file content is invalid JSON")
					}

					pricing = string(pricingContent)
				}

				buf := bytes.NewBuffer([]byte{})
				if err := json.Compact(buf, []byte(pricing)); err != nil {
					return fmt.Errorf("failed to compact the pricing")
				}

				pricing = buf.String()
			}

			qos := uint64(viper.GetInt64(FlagQoS))

			msg := types.NewMsgUpdateServiceBinding(args[0], provider, deposit, pricing, qos, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsUpdateServiceBinding)

	return cmd
}

// GetCmdSetWithdrawAddr implements setting a withdrawal address command
func GetCmdSetWithdrawAddr(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-withdraw-addr [withdrawal-address]",
		Short: "Set a withdrawal address for an owner",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Set a withdrawal address for an owner.

Example:
$ %s tx service set-withdraw-addr <withdrawal-address> --from mykey
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			withdrawAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetWithdrawAddress(owner, withdrawAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdDisableServiceBinding implements disabling a service binding command
func GetCmdDisableServiceBinding(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable [service-name] [provider-address]",
		Short: "Disable an available service binding",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Disable an available service binding.

Example:
$ %s tx service disable <service-name> <provider-address> --from mykey
`,
				version.ClientName,
			),
		),
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			var err error
			var provider sdk.AccAddress

			if len(args) > 1 {
				provider, err = sdk.AccAddressFromBech32(args[1])
				if err != nil {
					return err
				}
			} else {
				provider = owner
			}

			msg := types.NewMsgDisableServiceBinding(args[0], provider, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdEnableServiceBinding implements enabling a service binding command
func GetCmdEnableServiceBinding(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable [service-name] [provider-address]",
		Short: "Enable an unavailable service binding",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Enable an unavailable service binding.

Example:
$ %s tx service enable <service-name> <provider-address> --deposit=1stake --from mykey
`,
				version.ClientName,
			),
		),
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			var err error
			var provider sdk.AccAddress

			if len(args) > 1 {
				provider, err = sdk.AccAddressFromBech32(args[1])
				if err != nil {
					return err
				}
			} else {
				provider = owner
			}

			var deposit sdk.Coins

			depositStr := viper.GetString(FlagDeposit)
			if len(depositStr) != 0 {
				deposit, err = sdk.ParseCoins(depositStr)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgEnableServiceBinding(args[0], provider, deposit, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsEnableServiceBinding)

	return cmd
}

// GetCmdRefundServiceDeposit implements refunding deposit command
func GetCmdRefundServiceDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refund-deposit [service-name] [provider-address]",
		Short: "Refund all deposit from a service binding",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Refund all deposit from a service binding.

Example:
$ %s tx service refund-deposit <service-name> <provider-address> --from mykey
`,
				version.ClientName,
			),
		),
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			var err error
			var provider sdk.AccAddress

			if len(args) > 1 {
				provider, err = sdk.AccAddressFromBech32(args[1])
				if err != nil {
					return err
				}
			} else {
				provider = owner
			}

			msg := types.NewMsgRefundServiceDeposit(args[0], provider, owner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdCallService implements initiating a service call command
func GetCmdCallService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call",
		Short: "Initiate a service call",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Initiate a service call.

Example:
$ %s tx service call --service-name=<service-name> --providers=<provider list> 
--service-fee-cap=1stake --data=<input content or path/to/input.json> --timeout=100 
--repeated --frequency=150 --total=100 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			consumer := cliCtx.GetFromAddress()

			serviceName := viper.GetString(FlagServiceName)

			var providers []sdk.AccAddress
			providerList := viper.GetStringSlice(FlagProviders)

			for _, p := range providerList {
				provider, err := sdk.AccAddressFromBech32(p)
				if err != nil {
					return err
				}

				providers = append(providers, provider)
			}

			serviceFeeCap, err := sdk.ParseCoins(viper.GetString(FlagServiceFeeCap))
			if err != nil {
				return err
			}

			input := viper.GetString(FlagData)

			if !json.Valid([]byte(input)) {
				inputContent, err := ioutil.ReadFile(input)
				if err != nil {
					return fmt.Errorf("invalid input data: neither JSON input nor path to .json file were provided")
				}

				if !json.Valid(inputContent) {
					return fmt.Errorf("invalid input data: .json file content is invalid JSON")
				}

				input = string(inputContent)
			}

			buf := bytes.NewBuffer([]byte{})
			if err := json.Compact(buf, []byte(input)); err != nil {
				return fmt.Errorf("failed to compact the input data")
			}

			input = buf.String()
			timeout := viper.GetInt64(FlagTimeout)
			superMode := viper.GetBool(FlagSuperMode)
			repeated := viper.GetBool(FlagRepeated)

			frequency := uint64(0)
			total := int64(0)

			if repeated {
				frequency = uint64(viper.GetInt64(FlagFrequency))
				total = viper.GetInt64(FlagTotal)
			}

			msg := types.NewMsgCallService(
				serviceName, providers, consumer, input, serviceFeeCap,
				timeout, superMode, repeated, frequency, total,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsCallService)
	_ = cmd.MarkFlagRequired(FlagServiceName)
	_ = cmd.MarkFlagRequired(FlagProviders)
	_ = cmd.MarkFlagRequired(FlagServiceFeeCap)
	_ = cmd.MarkFlagRequired(FlagData)
	_ = cmd.MarkFlagRequired(FlagTimeout)

	return cmd
}

// GetCmdRespondService implements responding to a service request command
func GetCmdRespondService(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "respond",
		Short: "Respond to a service request",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Respond to an active service request.

Example:
$ %s tx service respond --request-id=<request-id> --result=<result content or path/to/result.json>
--data=<output content or path/to/output.json> --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			provider := cliCtx.GetFromAddress()

			requestIDStr := viper.GetString(FlagRequestID)
			requestID, err := types.ConvertRequestID(requestIDStr)
			if err != nil {
				return err
			}

			result := viper.GetString(FlagResult)
			output := viper.GetString(FlagData)

			if len(result) > 0 {
				if !json.Valid([]byte(result)) {
					resultContent, err := ioutil.ReadFile(result)
					if err != nil {
						return fmt.Errorf("invalid result: neither JSON input nor path to .json file were provided")
					}

					if !json.Valid(resultContent) {
						return fmt.Errorf("invalid result: .json file content is invalid JSON")
					}

					result = string(resultContent)
				}

				buf := bytes.NewBuffer([]byte{})
				if err := json.Compact(buf, []byte(result)); err != nil {
					return fmt.Errorf("failed to compact the result")
				}

				result = buf.String()
			}

			if len(output) > 0 {
				if !json.Valid([]byte(output)) {
					outputContent, err := ioutil.ReadFile(output)
					if err != nil {
						return fmt.Errorf("invalid output data: neither JSON input nor path to .json file were provided")
					}

					if !json.Valid(outputContent) {
						return fmt.Errorf("invalid output data: .json file content is invalid JSON")
					}

					output = string(outputContent)
				}

				buf := bytes.NewBuffer([]byte{})
				if err := json.Compact(buf, []byte(output)); err != nil {
					return fmt.Errorf("failed to compact the output data")
				}

				output = buf.String()
			}

			msg := types.NewMsgRespondService(requestID, provider, result, output)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsRespondService)
	_ = cmd.MarkFlagRequired(FlagRequestID)
	_ = cmd.MarkFlagRequired(FlagResult)

	return cmd
}

// GetCmdPauseRequestContext implements pausing a request context command
func GetCmdPauseRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause [request-context-id]",
		Short: "Pause a running request context",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Pause a running request context.

Example:
$ %s tx service pause <request-context-id> --from mykey
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			consumer := cliCtx.GetFromAddress()

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgPauseRequestContext(requestContextID, consumer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdStartRequestContext implements restarting a request context command
func GetCmdStartRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [request-context-id]",
		Short: "Start a paused request context",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Start a paused request context.

Example:
$ %s tx service start <request-context-id> --from mykey
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			consumer := cliCtx.GetFromAddress()

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgStartRequestContext(requestContextID, consumer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdKillRequestContext implements terminating a request context command
func GetCmdKillRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kill [request-context-id]",
		Short: "Terminate a request context",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Terminate a request context.

Example:
$ %s tx service kill <request-context-id> --from mykey
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			consumer := cliCtx.GetFromAddress()

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgKillRequestContext(requestContextID, consumer)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdUpdateRequestContext implements updating a request context command
func GetCmdUpdateRequestContext(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [request-context-id]",
		Short: "Update a request context",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Update a request context.

Example:
$ %s tx service update <request-context-id> --providers=<new providers> 
--service-fee-cap=2iris --timeout=0 --frequency=200 --total=200 --from mykey
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			consumer := cliCtx.GetFromAddress()

			requestContextID, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			var providers []sdk.AccAddress
			providerList := viper.GetStringSlice(FlagProviders)

			for _, p := range providerList {
				provider, err := sdk.AccAddressFromBech32(p)
				if err != nil {
					return err
				}

				providers = append(providers, provider)
			}

			var serviceFeeCap sdk.Coins

			serviceFeeCapStr := viper.GetString(FlagServiceFeeCap)
			if len(serviceFeeCapStr) != 0 {
				serviceFeeCap, err = sdk.ParseCoins(serviceFeeCapStr)
				if err != nil {
					return err
				}
			}

			timeout := viper.GetInt64(FlagTimeout)
			frequency := uint64(viper.GetInt64(FlagFrequency))
			total := viper.GetInt64(FlagTotal)

			msg := types.NewMsgUpdateRequestContext(
				requestContextID, providers, serviceFeeCap,
				timeout, frequency, total, consumer,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsUpdateRequestContext)

	return cmd
}

// GetCmdWithdrawEarnedFees implements withdrawing earned fees command
func GetCmdWithdrawEarnedFees(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-fees [provider-address]",
		Short: "Withdraw the earned fees of a provider or owner",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw the earned fees of the specified provider, 
for all providers of the owner if the provider not given.

Example:
$ %s tx service withdraw-fees <provider-address> --from mykey
`,
				version.ClientName,
			),
		),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			owner := cliCtx.GetFromAddress()

			var err error
			var provider sdk.AccAddress

			if len(args) > 1 {
				provider, err = sdk.AccAddressFromBech32(args[1])
				if err != nil {
					return err
				}
			} else {
				provider = owner
			}

			msg := types.NewMsgWithdrawEarnedFees(owner, provider)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
