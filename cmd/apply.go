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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/terrarium-tf/cli/lib"
)

func NewApplyCommand(root *cobra.Command) {
	var applyCmd = &cobra.Command{
		Use:   "apply workspace path/to/stack",
		Short: "Apply a given Terraform Stack",
		Long:  `Creates a plan file (which might be uploaded to CI-Artifacts for auditing) and applies this exact plan file.`,
		Args:  lib.ArgsValidator,

		Run: func(cmd *cobra.Command, args []string) {
			tf, ctx, files, _ := lib.Executor(*cmd, args[0], args[1], true)

			planFile := fmt.Sprintf("%s-%s.tfplan", strings.Replace(time.Now().Format(time.RFC3339), ":", "-", -1), args[0])
			planFile, _ = filepath.Abs(planFile)

			//plan
			_, _ = tf.Plan(ctx, buildPlanOptions(files, args, planFile)...)

			//apply
			_ = tf.Apply(ctx, tfexec.DirOrPlan(planFile))

			if os.Getenv("TF_IN_AUTOMATION") == "" {
				_ = os.Remove(planFile)
			}
		},
	}

	root.AddCommand(applyCmd)
}

func buildPlanOptions(files []string, args []string, planFile string) []tfexec.PlanOption {
	var planops []tfexec.PlanOption

	for _, f := range files {
		planops = append(planops, tfexec.VarFile(f))
	}

	if planFile != "" {
		return append(planops, tfexec.Var(fmt.Sprintf("environment=%s", args[0])), tfexec.Out(planFile))
	}

	return append(planops, tfexec.Var(fmt.Sprintf("environment=%s", args[0])))
}
