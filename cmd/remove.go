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
	"errors"
	"github.com/spf13/cobra"
	"github.com/terrarium-tf/cli/lib"
)

func NewRemoveCommand(root *cobra.Command) {
	var removeCmd = &cobra.Command{
		Use:   "remove workspace stack tf_resource_id",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: removeArgsValidator,
		Run: func(cmd *cobra.Command, args []string) {
			tf, ctx := lib.Executor(*cmd, args[0], args[1])
			_, _ = lib.Vars(*cmd, args[0], args[1])

			_ = tf.StateRm(ctx, args[2])
		},
	}

	root.AddCommand(removeCmd)
}

func removeArgsValidator(cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		return errors.New("requires a workspace,a stack path, a remote resource and a tf resource")
	}

	return lib.ArgsValidator(cmd, args)
}
