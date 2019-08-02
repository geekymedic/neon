/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	mini_gateway "github.com/geekymedic/neonx-cli/mini-gateway"
	"github.com/geekymedic/neonx-cli/types"
	"github.com/geekymedic/neonx-cli/util"
	"github.com/spf13/cobra"
	"os"
)

var generateServiceCmdOpt = struct {
	Name     string
	implName string
}{}

// serviceCmd represents the service command
var generateServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		// build system
		dir, err := util.AbsDir(dir)
		if err != nil {
			util.StdoutExit(-1, "Fail to generate service: %v", err)
		}
		sysDirNode := types.NewBaseDir(dir)
		err = mini_gateway.GenerateService(sysDirNode, generateServiceCmdOpt.Name)
		if err != nil {
			util.StdoutExit(-1, "Fail to create service: %v", err)
		}
		util.StdoutOk("Create service successfully\n")
	},
}

func init() {
	generateCmd.AddCommand(generateServiceCmd)
	generateServiceCmd.Flags().StringVar(&generateServiceCmdOpt.Name, "name", "CheckHealth", "service name")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}