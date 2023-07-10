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
	"os"
)

func NewRemoveCommand(root *cobra.Command) {
	var removeCmd = &cobra.Command{
		Use:     "remove workspace path/to/stack tf_resource_id",
		Short:   "Removes a remote resource from the terraform state",
		Example: "remove prod path/to/stack aws_s3_bucket.example",
		Args:    removeArgsValidator,
		Run: func(cmd *cobra.Command, args []string) {
			tf, ctx, _, _ := lib.Executor(*cmd, args[0], args[1], true)

			err := tf.StateRm(ctx, args[2])
			if err != nil {
				os.Exit(1)
			}
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
