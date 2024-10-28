// Copyright 2022 Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	log "github.com/inconshreveable/log15"
	"github.com/jeltevanbommel/scion-windows/environment"
	"github.com/jeltevanbommel/scion-windows/external/bootstrapper"
	"github.com/jeltevanbommel/scion-windows/external/daemon"
	"github.com/scionproto/scion/pkg/private/serrors"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
)

func newEndhost(pather CommandPather) *cobra.Command {
	var flags struct {
		bootstrapUrl string
	}

	var cmd = &cobra.Command{
		Use:     "run [flags]",
		Short:   "Bootstrap and run the SCION daemon for end host connectivity",
		Example: fmt.Sprintf("  %[1]s bootstrap", pather.CommandPath()),
		Long:    `'run' fetches the TRCs and topology files for your specific SCION AS using the bootstrapper and then starts the SCION daemon with the fetched configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if environment.EndhostEnv.Windows && !isAdmin() {
				runElevated()
				fmt.Println("Bootstrapping requires admin privileges. The application was run as a standard user. Rerunning with admin privileges...")
				return nil
			}
			endhostEnv := environment.EndhostEnv
			if flags.bootstrapUrl != "" && cmd.Flags().Lookup("url").Changed {
				endhostEnv.BootstrappingUrl = flags.bootstrapUrl
			}

			endhostEnv.Install()

			code := bootstrapper.Run(filepath.Join(endhostEnv.ConfigPath, "bootstrapper.toml"), endhostEnv.ConfigPath)
			if code != 0 {
				return serrors.New("Bootstrapping failed!")
			}
			log.Info("Running daemon...")
			return daemon.RunDaemon(context.Background(), path.Join(endhostEnv.DaemonConfigPath, "sciond.toml"))
		},
	}

	cmd.Flags().StringVar(&flags.bootstrapUrl, "url", "human",
		"Specify the url of the bootstrap server")

	return cmd
}
