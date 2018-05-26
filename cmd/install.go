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
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a windows service",
	Long: `This will install a windows service with the provided name and description.
	
	Example usage:

	svcmanager install -n "MyService" -d "My Service pulls data from the internet" -e "c//services/myservice.exe" -m

	For more information on how to use the command type

	svcmanager install --help 

	The user account to use for starting the service is optional. 
	`,
	RunE: runInstallCmd,
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&argServiceName, "name", "n", "", "Name of service")
	installCmd.Flags().StringVarP(&argServiceDescription, "description", "d", "", "Description of service")
	installCmd.Flags().StringVarP(&argExecutablePath, "path", "e", "", "Path to executable")
	installCmd.Flags().BoolVarP(&argManual, "manual", "m", false, "Don't start service on boot")
	installCmd.Flags().StringVarP(&argUser, "user", "u", "", "Username to start service as, instead of the LOCAL SYSTEM account. If no domain is given then the local system will be provided.")
	installCmd.Flags().StringVarP(&argPassword, "password", "p", "", "Password for the specified user. If not supplied then it is prompted.")
}

func runInstallCmd(cmd *cobra.Command, args []string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	svc, err := m.OpenService(argServiceName)
	if err == nil {
		svc.Close()
		return fmt.Errorf("service %s already exists", argServiceName)
	}

	cfg := mgr.Config{
		DisplayName: argServiceName,
		Description: argServiceDescription,
	}

	if argUser != "" {
		if strings.Index(argUser, "\\") < 0 && strings.Index(argUser, "@") < 0 {
			domain, err := windows.ComputerName()
			if err != nil || domain == "" {
				return errors.New("Unable to determine computer name; supply a fully qualified user instead")
			}
			argUser = domain + "\\" + argUser
			fmt.Fprintf(os.Stderr, "Changed user account to %s\n", argUser)
		}
		if argPassword == "" {
			fmt.Fprintf(os.Stderr, "Password for account %s: ", argUser)
			pwd, err := gopass.GetPasswd()
			if err != nil {
				return err
			}
			argPassword = string(pwd)
		}
		cfg.ServiceStartName = argUser
		cfg.Password = argPassword
	}

	if !argManual {
		cfg.StartType = mgr.StartAutomatic
	}

	svc, err = m.CreateService(argServiceName, argExecutablePath, cfg)
	if err != nil {
		return err
	}
	defer svc.Close()

	err = eventlog.InstallAsEventCreate(argServiceName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		svc.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}

	fmt.Fprintf(os.Stderr, "Installed service %s\n", argServiceName)
	return nil
}
