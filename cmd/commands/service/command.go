/*
 * Copyright (C) 2017 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mysteriumnetwork/node/cmd"
	"github.com/mysteriumnetwork/node/config"
	"github.com/mysteriumnetwork/node/config/urfavecli/clicontext"
	"github.com/mysteriumnetwork/node/core/node"
	"github.com/mysteriumnetwork/node/metadata"
	"github.com/mysteriumnetwork/node/services"
	"github.com/mysteriumnetwork/node/tequilapi/client"
	"github.com/mysteriumnetwork/node/tequilapi/contract"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

// NewCommand function creates service command
func NewCommand(licenseCommandName string) *cli.Command {
	var di cmd.Dependencies
	command := &cli.Command{
		Name:      "service",
		Usage:     "Starts and publishes services on Mysterium Network",
		ArgsUsage: "comma separated list of services to start",
		Before:    clicontext.LoadUserConfigQuietly,
		Action: func(ctx *cli.Context) error {
			if !ctx.Bool(config.FlagAgreedTermsConditions.Name) {
				printTermWarning(licenseCommandName)
				os.Exit(2)
			}

			quit := make(chan error)
			config.ParseFlagsServiceStart(ctx)
			config.ParseFlagsServiceOpenvpn(ctx)
			config.ParseFlagsServiceWireguard(ctx)
			config.ParseFlagsServiceNoop(ctx)
			config.ParseFlagsNode(ctx)

			nodeOptions := node.GetOptions()
			nodeOptions.Discovery.FetchEnabled = false
			if err := di.Bootstrap(*nodeOptions); err != nil {
				return err
			}
			go func() { quit <- di.Node.Wait() }()

			cmd.RegisterSignalCallback(func() { quit <- nil })

			cmdService := &serviceCommand{
				tequilapi:    client.NewClient(nodeOptions.TequilapiAddress, nodeOptions.TequilapiPort),
				errorChannel: quit,
			}
			go func() {
				quit <- cmdService.Run(ctx)
			}()

			return describeQuit(<-quit)
		},
		After: func(ctx *cli.Context) error {
			return di.Shutdown()
		},
	}

	config.RegisterFlagsServiceStart(&command.Flags)
	config.RegisterFlagsServiceOpenvpn(&command.Flags)
	config.RegisterFlagsServiceWireguard(&command.Flags)
	config.RegisterFlagsServiceNoop(&command.Flags)

	return command
}

func describeQuit(err error) error {
	if err == nil {
		log.Info().Msg("Stopping application")
	} else {
		log.Error().Err(err).Stack().Msg("Terminating application due to error")
	}
	return err
}

// serviceCommand represent entrypoint for service command with top level components
type serviceCommand struct {
	tequilapi    *client.Client
	errorChannel chan error
}

// Run runs a command
func (sc *serviceCommand) Run(ctx *cli.Context) (err error) {
	arg := ctx.Args().Get(0)
	serviceTypes := services.Types()
	if arg != "" {
		serviceTypes = strings.Split(arg, ",")
	}

	providerID := sc.unlockIdentity(
		ctx.String(config.FlagIdentity.Name),
		ctx.String(config.FlagIdentityPassphrase.Name),
	)
	log.Info().Msgf("Unlocked identity: %v", providerID)

	for _, serviceType := range serviceTypes {
		serviceOpts, err := services.GetStartOptions(serviceType)
		if err != nil {
			return err
		}
		startRequest := contract.ServiceStartRequest{
			ProviderID: providerID,
			Type:       serviceType,
			PaymentMethod: contract.ServicePaymentMethod{
				PriceGB:     serviceOpts.PaymentPricePerGB,
				PriceMinute: serviceOpts.PaymentPricePerMinute,
			},
			AccessPolicies: contract.ServiceAccessPolicies{IDs: serviceOpts.AccessPolicyList},
			Options:        serviceOpts,
		}

		go sc.runService(startRequest)
	}

	return <-sc.errorChannel
}

func (sc *serviceCommand) unlockIdentity(id, passphrase string) string {
	const retryRate = 10 * time.Second
	for {
		id, err := sc.tequilapi.CurrentIdentity(id, passphrase)
		if err == nil {
			return id.Address
		}
		log.Warn().Err(err).Msg("Failed to get current identity")
		log.Warn().Msgf("retrying in %vs...", retryRate.Seconds())
		time.Sleep(retryRate)
	}
}

func (sc *serviceCommand) runService(request contract.ServiceStartRequest) {
	_, err := sc.tequilapi.ServiceStart(request)
	if err != nil {
		sc.errorChannel <- errors.Wrapf(err, "failed to run service %s", request.Type)
	}
}

func printTermWarning(licenseCommandName string) {
	fmt.Println(metadata.VersionAsSummary(metadata.LicenseCopyright(
		"run program with 'myst "+licenseCommandName+" --"+config.LicenseWarrantyFlag.Name+"' option",
		"run program with 'myst "+licenseCommandName+" --"+config.LicenseConditionsFlag.Name+"' option",
	)))
	fmt.Println()

	fmt.Println("If you agree with these Terms & Conditions, run program again with '--agreed-terms-and-conditions' flag")
}
