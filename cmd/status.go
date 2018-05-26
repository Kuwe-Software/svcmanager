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

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Checks the status of a windows service",
	Long: `This will check the status of the windows service with the provided name.
	
	Example usage:

	svcmanager status --name "MyService" 

	or

	svcmanager status --n "MyService" 

	For more information on how to use the command type

	svcmanager install --help 
	`,
	RunE: runStatusCmd,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringVarP(&argServiceName, "name", "n", "", "Name of service")
}

func runStatusCmd(cmd *cobra.Command, args []string) error {
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

	status, err := s.Query()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s: %s\n", s.Name, stateNames[status.State])
	return nil
}
