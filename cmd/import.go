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
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
	"github.com/terrarium-tf/cli/lib"
)

func NewImportCommand(root *cobra.Command) {
	var importCmd = &cobra.Command{
		Use:     "import workspace path/to/stack tf_resource_id remote_resource",
		Short:   "Import a remote resource into a local terraform resource",
		Example: "import prod path/to/stack aws_s3_bucket.example some_aws_bucket_name",
		Args:    importArgsValidator,
		RunE: func(cmd *cobra.Command, args []string) error {
			tf, ctx, files, _ := lib.Executor(*cmd, args[0], args[1], true)

			return tf.Import(ctx, args[2], args[3], buildImportOptions(files, args)...)
		},
	}

	root.AddCommand(importCmd)
}

func importArgsValidator(cmd *cobra.Command, args []string) error {
	if len(args) < 4 {
		return errors.New("requires a workspace,a stack path, a remote resource and a tf resource")
	}

	return lib.ArgsValidator(cmd, args)
}

func buildImportOptions(files []string, args []string) []tfexec.ImportOption {
	var ops []tfexec.ImportOption

	for _, f := range files {
		ops = append(ops, tfexec.VarFile(f))
	}

	return append(ops, tfexec.Var(fmt.Sprintf("environment=%s", args[0])))
}
