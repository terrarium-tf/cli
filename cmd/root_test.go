package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	_, err = root.ExecuteC()

	return buf.String(), err
}

func runCommand(t *testing.T, args []string) string {
	rc := NewRootCommand()
	AddChildCommands(rc)
	output, err := executeCommand(rc, args...)
	if output == "" {
		t.Errorf("Unexpected output: %v", output)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	return output
}

func _varFilesArgs(t *testing.T) string {
	global, err := filepath.Abs("./../example/global.tfvars.json")
	if err != nil {
		t.Error(err)
	}
	app, err := filepath.Abs("./../example/stack/app.tfvars.json")
	if err != nil {
		t.Error(err)
	}
	dev, err := filepath.Abs("./../example/stack/dev.tfvars.json")
	if err != nil {
		t.Error(err)
	}

	return fmt.Sprintf("-var-file=%s -var-file=%s -var-file=%s", global, app, dev)
}

func TestInitCommandWithoutRemoteState(t *testing.T) {
	args := []string{"init", "dev", "../example/stack", "-t", "echo", "--remote-state=false"}
	out := runCommand(t, args)
	t.Log(out)

	if !strings.Contains(out, "init -force-copy -input=false -backend=false -get=true -upgrade=true") {
		t.Errorf("invalid init command")
	}
}

func TestInitCommandWithRemoteStateButNoLocking(t *testing.T) {
	args := []string{"init", "dev", "../example/stack", "-t", "echo", "--state-lock=false"}
	out := runCommand(t, args)
	t.Log(out)

	if !strings.Contains(out, "init -force-copy -input=false -backend=true -get=true -upgrade=true -backend-config=region=eu-central-1 -backend-config=bucket=tf-state-terrarium-cli-eu-central-1-455201159890 -backend-config=key=stack.tfstate") {
		t.Errorf("invalid init command")
	}
}

func TestInitCommand(t *testing.T) {
	args := []string{"init", "dev", "../example/stack", "-t", "echo"}
	out := runCommand(t, args)
	t.Log(out)

	if !strings.Contains(out, "init -force-copy -input=false -backend=true -get=true -upgrade=true -backend-config=region=eu-central-1 -backend-config=bucket=tf-state-terrarium-cli-eu-central-1-455201159890 -backend-config=key=stack.tfstate -backend-config=dynamodb_table=terraform-lock-terrarium-cli-eu-central-1-455201159890") {
		t.Errorf("invalid init command")
	}
}

func TestInitCommandGcp(t *testing.T) {
	args := []string{"init", "dev", "../example/stack_gcp", "-t", "echo"}
	out := runCommand(t, args)
	t.Log(out)

	if !strings.Contains(out, "init -force-copy -input=false -backend=true -get=true -upgrade=true -backend-config=credentials=foo -backend-config=bucket=tf-state-terrarium-cli -backend-config=prefix=stack_gcp") {
		t.Errorf("invalid init command")
	}
}

func TestInitCommandAzure(t *testing.T) {
	args := []string{"init", "dev", "../example/stack_azure", "-t", "echo"}
	out := runCommand(t, args)
	t.Log(out)

	if !strings.Contains(out, "init -force-copy -input=false -backend=true -get=true -upgrade=true -backend-config=storage_account_name=terrariumaccount -backend-config=resource_group_name=terrarium -backend-config=key=terrarium.tfstate -backend-config=container_name=tf-state-terrarium-cli-terrariumaccount") {
		t.Errorf("invalid init command")
	}
}

func TestApplyCommand(t *testing.T) {
	args := []string{"apply", "dev", "../example/stack", "-t", "echo"}
	now := strings.Replace(time.Now().Format(time.RFC3339), ":", "-", -1)
	out := runCommand(t, args)
	t.Log(out)

	root, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	if !strings.Contains(out, "workspace new dev") {
		t.Errorf("missing create workspace")
	}
	if !strings.Contains(out, "workspace select dev") {
		t.Errorf("missing switch workspace")
	}
	if !strings.Contains(out, fmt.Sprintf("plan -input=false -detailed-exitcode -lock-timeout=0s -out=%s/%s-dev.tfplan %s -lock=true -parallelism=10 -refresh=true -var environment=dev", root, now, _varFilesArgs(t))) {
		t.Errorf("invalid plan command")
	}
	if !strings.Contains(out, fmt.Sprintf("apply -auto-approve -input=false -lock=true -parallelism=10 -refresh=true %s/%s-dev.tfplan", root, now)) {
		t.Errorf("invalid apply command")
	}
}

func TestPlanCommand(t *testing.T) {
	args := []string{"plan", "dev", "../example/stack", "-t", "echo"}
	out := runCommand(t, args)
	t.Log(out)

	if !strings.Contains(out, "workspace new dev") {
		t.Errorf("missing create workspace")
	}
	if !strings.Contains(out, "workspace select dev") {
		t.Errorf("missing switch workspace")
	}
	if !strings.Contains(out, fmt.Sprintf("plan -input=false -detailed-exitcode -lock-timeout=0s %s -lock=true -parallelism=10 -refresh=true -var environment=dev", _varFilesArgs(t))) {
		t.Errorf("invalid plan command")
	}
}

func TestDestroyCommand(t *testing.T) {
	args := []string{"destroy", "dev", "../example/stack", "-t", "echo"}
	out := runCommand(t, args)
	t.Log(out)

	if !strings.Contains(out, "workspace new dev") {
		t.Errorf("missing create workspace")
	}
	if !strings.Contains(out, "workspace select dev") {
		t.Errorf("missing switch workspace")
	}
	if !strings.Contains(out, fmt.Sprintf("destroy -auto-approve -input=false -lock-timeout=0s %s -lock=true -parallelism=10 -refresh=true -var environment=dev", _varFilesArgs(t))) {
		t.Errorf("invalid destroy command")
	}
}

func TestImportCommand(t *testing.T) {
	args := []string{"import", "dev", "../example/stack", "-t", "echo", "aws_s3_bucket.test", "some_bucket_id"}
	out := runCommand(t, args)
	t.Log(out)

	if !strings.Contains(out, "workspace new dev") {
		t.Errorf("missing create workspace")
	}
	if !strings.Contains(out, "workspace select dev") {
		t.Errorf("missing switch workspace")
	}
	if !strings.Contains(out, fmt.Sprintf("import -input=false -lock-timeout=0s %s -lock=true -var environment=dev aws_s3_bucket.test some_bucket_id", _varFilesArgs(t))) {
		t.Errorf("invalid import command")
	}
}

func TestRemoveCommand(t *testing.T) {
	args := []string{"remove", "dev", "../example/stack", "-t", "echo", "aws_s3_bucket.test"}
	out := runCommand(t, args)
	t.Log(out)

	if !strings.Contains(out, "workspace new dev") {
		t.Errorf("missing create workspace")
	}
	if !strings.Contains(out, "workspace select dev") {
		t.Errorf("missing switch workspace")
	}
	if !strings.Contains(out, "state rm -lock-timeout=0s -lock=true aws_s3_bucket.test") {
		t.Errorf("invalid remove command")
	}
}

func TestRemoveCommandWithVerbose(t *testing.T) {
	args := []string{"remove", "dev", "../example/stack", "-t", "echo", "aws_s3_bucket.test", "-v"}
	out := runCommand(t, args)
	t.Log(out)

	if !strings.Contains(out, "workspace new dev") {
		t.Errorf("missing create workspace")
	}
	if !strings.Contains(out, "workspace select dev") {
		t.Errorf("missing switch workspace")
	}
	if !strings.Contains(out, "Collected vars files:") {
		t.Errorf("missing verbose log about var files")
	}
	if !strings.Contains(out, "Collected vars:") {
		t.Errorf("missing verbose log about vars")
	}
	if !strings.Contains(out, "state rm -lock-timeout=0s -lock=true aws_s3_bucket.test") {
		t.Errorf("invalid remove command")
	}
}

func TestUntaintCommand(t *testing.T) {
	t.Skip("test not yet fully working")
	args := []string{"untaint", "dev", "../example/stack", "-t", "echo", "aws_s3_bucket.test"}
	out := runCommand(t, args)
	t.Log(out)

	if !strings.Contains(out, "workspace new dev") {
		t.Errorf("missing create workspace")
	}
	if !strings.Contains(out, "workspace select dev") {
		t.Errorf("missing switch workspace")
	}
	if !strings.Contains(out, "untaint aws_s3_bucket.test") {
		t.Errorf("invalid untaint command")
	}
}
