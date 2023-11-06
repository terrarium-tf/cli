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
	"github.com/spf13/cobra"
	"github.com/terrarium-tf/cli/lib"
	"os"
	"strings"
	"time"
)

func NewPlanCommand(root *cobra.Command) {
	var planCmd = &cobra.Command{
		Use:   "plan workspace path/to/stack",
		Short: "Creates a diff between remote and local state and prints the upcoming changes",
		Args:  lib.ArgsValidator,

		RunE: func(cmd *cobra.Command, args []string) error {
			tf, ctx, files, _ := lib.Executor(*cmd, args[0], args[1], true)

			//plan
			planFile := ""
			if os.Getenv("TF_IN_AUTOMATION") != "" {
				planFile = fmt.Sprintf("%s-%s.tfplan", strings.Replace(time.Now().Format(time.RFC3339), ":", "-", -1), args[0])
			}

			diff, err := tf.Plan(ctx, buildPlanOptions(files, args, planFile)...)

			// behave exactly like terraform:
			/*
				0 = Succeeded with empty diff (no changes)
				1 = Error
				2 = Succeeded with non-empty diff (changes present)
			*/
			if diff {
				os.Exit(2)
			}

			return err
		},
	}

	root.AddCommand(planCmd)
}
