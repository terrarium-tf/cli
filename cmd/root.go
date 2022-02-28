// Package cmd
/*
Copyright © 2022 Robert Schönthal <robert@schoenthal.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/terrarium-tf/cli/lib"
	"os"
)

func Execute(command *cobra.Command) {
	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	var binary string
	var verbose bool

	var rootCmd = &cobra.Command{
		Use:   "terrarium [command] workspace path/to/stack",
		Short: "run terraform with autoinjected var files",
		Long: `Builds Terraform Commands, easing those steps:
* collects defined var-files
* switches to the given workspace (can create new one)
* runs the given terraform command with the multiple -var-files options in correct order.

You can override the default terraform binary with "-t"
Add "-v" for more verbose logging.
`,
		Example: "terrarium [command] workspace path/to/stack -v -t echo",
	}

	rootCmd.PersistentFlags().StringVarP(&binary, "terraform", "t", lib.Binary(), "terraform binary found in your path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "display extended informations")

	return rootCmd
}

func AddChildCommands(rootCmd *cobra.Command) {
	NewApplyCommand(rootCmd)
	NewDestroyCommand(rootCmd)
	NewImportCommand(rootCmd)
	NewInitCommand(rootCmd)
	NewPlanCommand(rootCmd)
	NewRemoveCommand(rootCmd)
}
