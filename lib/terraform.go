package lib

import (
	"context"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

const (
	InfoColorLine    = "\033[1;34m%s\033[0m\n"
	NoticeColorLine  = "\033[1;36m%s\033[0m\n"
	WarningColorLine = "\033[1;33m%s\033[0m\n"
	WarningColorMap  = "\033[1;33m%-*s\033[0m : \u001B[1;33m%s\u001B[0m\n"
	ErrorColorLine   = "\033[1;31m%s\033[0m\n"
	DebugColorLine   = "\033[0;36m%s\033[0m\n"
)

func Binary() string {
	path, err := exec.LookPath("terraform")
	if err != nil {
		log.Fatal("cant find terraform binary, please provide it yourself", err)
	}
	return path
}

func Executor(cmd cobra.Command, workspace string, path string, switchWorkspace bool) (*tfexec.Terraform, context.Context, []string, map[string]any) {
	binary, err := cmd.Parent().PersistentFlags().GetString("terraform")

	if err != nil {
		cmd.PrintErr("cant find terraform flag", err)
	}

	tf, err := tfexec.NewTerraform(path, binary)
	tf.SetColor(true)

	if err != nil {
		cmd.PrintErr(err.Error())
	}

	tf.SetStdout(cmd.OutOrStdout())
	tf.SetStderr(cmd.ErrOrStderr())

	ctx := context.Background()

	if switchWorkspace {
		Workspace(tf, ctx, cmd, workspace)
	}

	files, vars := Vars(cmd, workspace, path)

	return tf, context.Background(), files, vars
}

func Workspace(tf *tfexec.Terraform, ctx context.Context, cmd cobra.Command, name string) {
	tf.SetStdout(nil)
	workspaces, current, err := tf.WorkspaceList(ctx)
	tf.SetStdout(cmd.OutOrStdout())

	if err != nil {
		cmd.PrintErr(err.Error())
	}

	exists := false
	for _, ws := range workspaces {
		if ws == name {
			exists = true
		}
	}
	if !exists {
		err := tf.WorkspaceNew(ctx, name)
		if err != nil {
			cmd.PrintErr(err.Error())
		}
	}

	if current != name {
		err := tf.WorkspaceSelect(ctx, name)
		if err != nil {
			cmd.PrintErr(err.Error())
		}
	}
}
