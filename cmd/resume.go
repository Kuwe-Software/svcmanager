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
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/spf13/cobra"
)

// resumeCmd represents the resume command
var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a windows service",
	Long: `This will resume the windows service with the provided name, if it is in a paused state.

	svcmanager resume --name "MyService" 

	or

	svcmanager resume --n "MyService" 

	For more information on how to use the command type

	svcmanager install --help 
	`,
	RunE: runResumeCmd,
}

func init() {
	rootCmd.AddCommand(resumeCmd)
	resumeCmd.Flags().StringVarP(&argServiceName, "name", "n", "", "Name of service")
}

func runResumeCmd(cmd *cobra.Command, args []string) error {
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

	status, err := s.Control(svc.Continue)
	if err != nil {
		return fmt.Errorf("could not resume service: %v", err)
	}

	timeout := time.Now().Add(10 * time.Second)
	for status.State != svc.Running {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", svc.Running)
		}
		time.Sleep(300 * time.Millisecond)

		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}

	fmt.Fprintf(os.Stderr, "Resumed service %s\n", argServiceName)
	return nil
}
