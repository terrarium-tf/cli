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
	"bufio"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/terrarium-tf/cli/lib"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func NewInitCommand(root *cobra.Command) {
	var initCmd = &cobra.Command{
		Use:   "init workspace path/to/stack [--remote-state=false] [--state-lock=false]",
		Short: "initializes a stack with optional remote state",
		Long: `The init command can (defaults to yes) configure the stack with a remote state.
All you need is

terraform {
  backend "s3" {
  }
}

The rest will be autogenerated.
Pattern for the bucket is: "tf-state-{PROJECT}-{REGION}-{ACCOUNT_ID}"
Pattern for the dynamo is: "terraform-lock-{PROJECT}-{REGION}-{ACCOUNT_ID}"

These variables can be defined by your *.tfvars.json or through command options
`,
		Example: "init workspace path/to/stack --state-bucket=my_own_bucket_id --state-dynamo=my_dynamo_table --state-region=us-east-1 --state-account=4711 --state-name=my_state_entry_name",
		Args:    lib.ArgsValidator,
		RunE: func(cmd *cobra.Command, args []string) error {
			tf, ctx, _, mergedVars := lib.Executor(*cmd, args[0], args[1], false)

			return tf.Init(ctx, buildInitOptions(*cmd, mergedVars, args)...)
		},
	}

	initCmd.Flags().BoolP("remote-state", "r", true, "initialize with remote state")

	// TODO allow more customizations for remote state
	initCmd.Flags().Bool("state-lock", true, "initialize with state locking")
	initCmd.Flags().String("state-bucket", "", "initialize with state bucket")
	initCmd.Flags().String("state-dynamo", "", "initialize with state dynamo for locking")
	initCmd.Flags().String("state-region", "", "initialize with state region")
	initCmd.Flags().String("state-account", "", "initialize with state aws|azure account")
	initCmd.Flags().String("state-name", "", "initialize with state name")

	root.AddCommand(initCmd)
}

func buildInitOptions(cmd cobra.Command, mergedVars map[string]any, args []string) []tfexec.InitOption {
	var opts []tfexec.InitOption

	// if we want to init with an s3/dynamo remote state
	rs := cmd.Flags().Lookup("remote-state")
	if rs.Value.String() == "true" {
		// find the backend provider by scanning files for a backend config statement
		switch detectBackendProvider(args[1]) {
		case "gcs":
			opts = configureGcsBackend(cmd, mergedVars, args, opts)
		case "azure":
			opts = configureAzureBackend(cmd, mergedVars, args, opts)
		default:
			opts = configureAwsBackend(cmd, mergedVars, args, opts)

		}
	} else {
		opts = append(opts, tfexec.Backend(false))
	}

	return append(opts, tfexec.Upgrade(true))
}

func detectBackendProvider(stackPath string) string {
	tfFiles := findFiles(stackPath, ".tf")

	for _, f := range tfFiles {
		provider, _ := scanFile(f)
		if provider != "" {
			return provider
		}
	}
	return "s3"
}

func scanFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "backend \"s3\"") {
			return "s3", nil
		}
		if strings.Contains(scanner.Text(), "backend \"gcs\"") {
			return "gcs", nil
		}
		if strings.Contains(scanner.Text(), "backend \"azurerm\"") {
			return "azure", nil
		}

		line++
	}

	if err = scanner.Err(); err != nil {
		// Handle the error
		return "", err
	}

	return "", nil
}

func findFiles(root, ext string) []string {
	var a []string
	_ = filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})

	return a
}

func configureAwsBackend(cmd cobra.Command, mergedVars map[string]any, args []string, opts []tfexec.InitOption) []tfexec.InitOption {
	opts = append(opts,
		configureRegion(cmd, mergedVars),
		configureAwsBucket(cmd, mergedVars),
		configureStateKey(cmd, mergedVars, args),
	)

	sl := cmd.Flags().Lookup("state-lock")
	if sl.Value.String() == "true" {
		opts = append(opts,
			configureStateLock(cmd, mergedVars),
		)
	}
	return opts
}

func configureGcsBackend(cmd cobra.Command, mergedVars map[string]any, args []string, opts []tfexec.InitOption) []tfexec.InitOption {
	return append(opts,
		configureGcpCredentials(cmd, mergedVars),
		configureGcpBucket(cmd, mergedVars),
		configurePrefix(cmd, mergedVars, args),
	)
}

func configureAzureBackend(cmd cobra.Command, mergedVars map[string]any, args []string, opts []tfexec.InitOption) []tfexec.InitOption {
	opts = append(opts,
		configureAzureCredentials(cmd, mergedVars),
		configureAzureGroupName(cmd, mergedVars),
		configureStateKey(cmd, mergedVars, args),
		configureAzureBucket(cmd, mergedVars),
	)

	return append(opts, configureAzureFromEnv(cmd, mergedVars)...)
}

func configureRegion(cmd cobra.Command, mergedVars map[string]any) tfexec.InitOption {
	region := lib.GetVar("region", cmd, mergedVars, false)

	if region == "" {
		if os.Getenv("AWS_REGION") != "" {
			region = os.Getenv("AWS_REGION")
		}
		if region == "" && os.Getenv("AWS_DEFAULT_REGION") != "" {
			region = os.Getenv("AWS_DEFAULT_REGION")
		}
	}

	if region == "" {
		log.Fatalf(lib.ErrorColorLine, "unable to configure remote state, 'region' was not found in var files and not provided with '-state-region' nor was AWS_REGION or AWS_DEFAULT_REGION found in global environment")
	}

	return tfexec.BackendConfig(fmt.Sprintf("region=%s", region))
}

func configureAzureCredentials(cmd cobra.Command, mergedVars map[string]any) tfexec.InitOption {
	account := lib.GetVar("account", cmd, mergedVars, true)

	return tfexec.BackendConfig(fmt.Sprintf("storage_account_name=%s", account))
}

func configureAzureFromEnv(cmd cobra.Command, mergedVars map[string]any) []tfexec.InitOption {
	var opts []tfexec.InitOption

	vars := [][]string{
		{"environment", "ARM_ENVIRONMENT"},
		{"endpoint", "ARM_ENDPOINT"},
		{"metadata_host", "ARM_METADATA_HOSTNAME"},
		{"snapshot", "ARM_SNAPSHOT"},
		{"msi_endpoint", "ARM_MSI_ENDPOINT"},
		{"use_msi", "ARM_USE_MSI"},
		{"oidc_request_url", "ARM_OIDC_REQUEST_URL"},
		{"oidc_request_token", "ARM_OIDC_REQUEST_TOKEN"},
		{"oidc_token", "ARM_OIDC_TOKEN"},
		{"oidc_token_file_path", "ARM_OIDC_TOKEN_FILE_PATH"},
		{"use_oidc", "ARM_USE_OIDC"},
		{"sas_token", "ARM_SAS_TOKEN"},
		{"access_key", "ARM_ACCESS_KEY"},
		{"use_azuread_auth", "ARM_USE_AZUREAD"},
		{"client_id", "ARM_CLIENT_ID"},
		{"client_certificate_password", "ARM_CLIENT_CERTIFICATE_PASSWORD"},
		{"client_certificate_path", "ARM_CLIENT_CERTIFICATE_PATH"},
		{"client_secret", "ARM_CLIENT_SECRET"},
		{"subscription_id", "ARM_SUBSCRIPTION_ID"},
		{"tenant_id", "ARM_TENANT_ID"},
	}

	for _, tuple := range vars {
		v := sourceAzureVars(tuple[0], tuple[1], cmd, mergedVars)
		if v != nil {
			opts = append(opts, v)
		}
	}

	return opts
}

func sourceAzureVars(varName string, envName string, cmd cobra.Command, mergedVars map[string]any) *tfexec.BackendConfigOption {
	v := lib.GetVar(varName, cmd, mergedVars, false)
	if v == "" && os.Getenv(envName) != "" {
		v = os.Getenv(envName)
	}
	if v != "" {
		return tfexec.BackendConfig(fmt.Sprintf("%s=%s", varName, v))
	}
	return nil
}

func configureAzureGroupName(cmd cobra.Command, mergedVars map[string]any) tfexec.InitOption {
	name := lib.GetVar("project", cmd, mergedVars, true)

	return tfexec.BackendConfig(fmt.Sprintf("resource_group_name=%s", name))
}

func configureGcpCredentials(cmd cobra.Command, mergedVars map[string]any) tfexec.InitOption {
	file := lib.GetVar("credentials", cmd, mergedVars, false)

	if file == "" {
		if os.Getenv("GOOGLE_BACKEND_CREDENTIALS") != "" {
			file = os.Getenv("GOOGLE_BACKEND_CREDENTIALS")
		}
		if file == "" && os.Getenv("GOOGLE_CREDENTIALS") != "" {
			file = os.Getenv("GOOGLE_CREDENTIALS")
		}
	}

	if file == "" {
		log.Fatalf(lib.ErrorColorLine, "unable to configure remote state, 'credentials' was not found in var files and not provided with '-state-credentials' nor was GOOGLE_BACKEND_CREDENTIALS or GOOGLE_CREDENTIALS found in global environment")
	}

	return tfexec.BackendConfig(fmt.Sprintf("credentials=%s", file))
}

func configureAwsBucket(cmd cobra.Command, mergedVars map[string]any) tfexec.InitOption {
	bucket := lib.GetVar("bucket", cmd, mergedVars, false)
	if bucket == "" {
		// no bucket defined, so generate a unique name
		bucket = fmt.Sprintf("tf-state-%s-%s-%s",
			lib.GetVar("project", cmd, mergedVars, false),
			lib.GetVar("region", cmd, mergedVars, false),
			lib.GetVar("account", cmd, mergedVars, true),
		)
	}
	return tfexec.BackendConfig(fmt.Sprintf("bucket=%s", bucket))
}

func configureAzureBucket(cmd cobra.Command, mergedVars map[string]any) tfexec.InitOption {
	bucket := lib.GetVar("bucket", cmd, mergedVars, false)
	if bucket == "" {
		// no bucket defined, so generate a unique name
		bucket = fmt.Sprintf("tf-state-%s-%s",
			lib.GetVar("project", cmd, mergedVars, false),
			lib.GetVar("account", cmd, mergedVars, true),
		)
	}
	return tfexec.BackendConfig(fmt.Sprintf("container_name=%s", bucket))
}

func configureGcpBucket(cmd cobra.Command, mergedVars map[string]any) tfexec.InitOption {
	bucket := lib.GetVar("bucket", cmd, mergedVars, false)
	if bucket == "" {
		// no bucket defined, so generate a unique name
		bucket = fmt.Sprintf("tf-state-%s",
			lib.GetVar("project", cmd, mergedVars, false),
		)
	}
	return tfexec.BackendConfig(fmt.Sprintf("bucket=%s", bucket))
}

func configureStateKey(cmd cobra.Command, mergedVars map[string]any, args []string) tfexec.InitOption {
	key := lib.GetVar("name", cmd, mergedVars, false)
	if key == "" {
		// no bucket defined, so generate a unique name
		key = path.Base(args[1])
	}
	return tfexec.BackendConfig(fmt.Sprintf("key=%s.tfstate", key))
}

func configurePrefix(cmd cobra.Command, mergedVars map[string]any, args []string) tfexec.InitOption {
	key := lib.GetVar("prefix", cmd, mergedVars, false)
	if key == "" {
		// no prefix defined, so generate a unique name
		key = path.Base(args[1])
	}
	return tfexec.BackendConfig(fmt.Sprintf("prefix=%s", key))
}

func configureStateLock(cmd cobra.Command, mergedVars map[string]any) tfexec.InitOption {
	table := lib.GetVar("dynamo", cmd, mergedVars, false)
	if table == "" {
		// no bucket defined, so generate a unique name
		table = fmt.Sprintf("terraform-lock-%s-%s-%s",
			lib.GetVar("project", cmd, mergedVars, false),
			lib.GetVar("region", cmd, mergedVars, false),
			lib.GetVar("account", cmd, mergedVars, true),
		)
	}
	return tfexec.BackendConfig(fmt.Sprintf("dynamodb_table=%s", table))
}
