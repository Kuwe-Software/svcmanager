// Copyright © 2018 Stewart Mbofana stewart.mbofana@live.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/svc/mgr"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a windows service",
	Long: `This will start the windows service with the provided name if it is not in a running state.

	Example usage:

	svcmanager start --name "MyService" 

	or

	svcmanager start --n "MyService" 

	For more information on how to use the command type

	svcmanager install --help 
	`,
	RunE: runStartCmd,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&argServiceName, "name", "n", "", "Name of service")
}

func runStartCmd(cmd *cobra.Command, args []string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(argServiceName)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}

	defer s.Close()

	err = s.Start("is", "manual-started")
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Started service %s\n", argServiceName)
	return nil
}
