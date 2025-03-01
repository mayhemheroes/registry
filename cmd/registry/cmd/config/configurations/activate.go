// Copyright 2022 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configurations

import (
	"fmt"

	"github.com/apigee/registry/pkg/config"
	"github.com/spf13/cobra"
)

func activateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate CONFIGURATION_NAME",
		Short: "Activates an existing named configuration.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			name := args[0]
			if err := config.ValidateName(name); err != nil {
				return err
			}

			_, err := config.Read(name)
			if err != nil {
				return fmt.Errorf("Cannot read config %q: %v", name, err)
			}

			err = config.Activate(name)
			if err != nil {
				return fmt.Errorf("Cannot activate config %q: %v", name, err)
			}

			cmd.Printf("Activated %q.\n", name)
			return nil
		},
	}
	return cmd
}
