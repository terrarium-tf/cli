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
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
	"github.com/terrarium-tf/cli/lib"
)

func NewDestroyCommand(root *cobra.Command) {
	var destroyCmd = &cobra.Command{
		Use:   "destroy workspace stack",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Args: lib.ArgsValidator,

		Run: func(cmd *cobra.Command, args []string) {
			tf, ctx := lib.Executor(*cmd, args[0], args[1])
			files, _ := lib.Vars(*cmd, args[0], args[1])

			_ = tf.Destroy(ctx, buildDestroyOptions(files, args)...)
		},
	}

	root.AddCommand(destroyCmd)
}

func buildDestroyOptions(files []string, args []string) []tfexec.DestroyOption {
	var ops []tfexec.DestroyOption

	for _, f := range files {
		ops = append(ops, tfexec.VarFile(f))
	}

	return append(ops, tfexec.Var(fmt.Sprintf("environment=%s", args[0])))
}
