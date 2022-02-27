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
	"time"

	"github.com/terrarium-tf/cli/lib"
)

func NewApplyCommand(root *cobra.Command) {
	var applyCmd = &cobra.Command{
		Use:   "apply workspace stack",
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
			planFile := fmt.Sprintf("%s-%s.tfplan", time.Now().Format(time.RFC3339), args[0])

			//plan
			_, _ = tf.Plan(ctx, buildPlanOptions(files, args, planFile)...)

			//apply
			_ = tf.Apply(ctx, tfexec.DirOrPlan(planFile))
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
