# Terrarium

> a tiny wrapper for Terraform to make loading env vars transparent by convention

[![Coverage Status](https://coveralls.io/repos/github/terrarium-tf/cli/badge.svg?branch=main)](https://coveralls.io/github/terrarium-tf/cli?branch=main)
[![Test & Build](https://github.com/terrarium-tf/cli/actions/workflows/test.yml/badge.svg)](https://github.com/terrarium-tf/cli/actions/workflows/test.yml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/terrarium-tf/cli)

**Builds Terraform Commands, easing these steps:**

* collects defined var-files
* switches to the given workspace (can create new one)
* runs the given terraform command with the multiple -var-files options in correct order.
* automatically detects `s3`, `gcs` or `azure` backend
* local file for machine only parameters

using these awesome tools:

* [terraform-exec](https://github.com/hashicorp/terraform-exec)
* [cobra](https://github.com/spf13/cobra)

## Setup

download the binary matching your OS from [here](https://github.com/terrarium-tf/cli/releases)

## Anatomy of Stacks and Configs

```
 |- local.tfvars.json # private variables available to all stacks, e.g. local paths (relative to cwd)
 |- global.tfvars.json # variables available to all stacks (relative to cwd)
 |- stage.tfvars.json # variables available to all stacks using the "stage" workspace (relative to cwd)
 |
 | - stacks
    |
    | - foo
        | - app.tfvars.json # default variables for this stack
        | - stage.tfvars.json # variables for the stage workspace (environment)
        |
        | - main.tf your stack entrypoint
```

## Command

```
$ terrarium
Builds Terraform Commands, easing these steps:
* collects defined var-files
* switches to the given workspace (can create new one)
* runs the given terraform command with the multiple -var-files options in correct order.

You can override the default terraform binary with "-t"
Add "-v" for more verbose logging.

Usage:
  terrarium [command]

Examples:
terrarium plan production path/to/stack -v

Available Commands:
  apply       Apply a given Terraform Stack
  completion  Generate the autocompletion script for the specified shell
  destroy     Destroy a given Terraform stack
  help        Help about any command
  import      Import a remote resource into a local terraform resource
  init        initializes a stack with optional remote state
  plan        Creates a diff between remote and local state and prints the upcoming changes
  remove      Removes a remote resource from the terraform state
  taint       Taints a given Terraform Resource from a State
  untaint     Untaints a given Terraform Resource from a State

Flags:
  -h, --help               help for terrarium
  -t, --terraform string   terraform binary found in your path (default "/usr/local/bin/terraform")
  -v, --verbose            display extended informations

Use "terrarium [command] --help" for more information about a command.
```

## Usage & Under the Hood

> assuming the above stack setup,

`terrarium init stage example/stack`

will internally run:

```
terraform workspace select stage
terraform version
terraform init -force-copy -input=false -backend=true -get=true -upgrade=true -backend-config=region=eu-central-1 -backend-config=bucket=tf-state-terrarium-cli-eu-central-1-455201159890 -backend-config=key=stack.tfstate -backend-config=dynamodb_table=terraform-lock-terrarium-cli-eu-central-1-455201159890
```

`terrarium apply stage example/stack`

will internally run:

```
terraform workspace select stage
terraform plan -input=false -detailed-exitcode -lock-timeout=0s -out=2022-02-28T16:26:26Z-stage.tfplan -var-file=example/global.tfvars.json -var-file=example/stack/app.tfvars.json -var-file=example/stack/stage.tfvars.json -lock=true -parallelism=10 -refresh=true -var environment=stage
terraform apply -auto-approve -input=false -lock=true -parallelism=10 -refresh=true 2022-02-28T16:26:26Z-stage.tfplan
```

## Usage in CI Runners

### Github-Actions

```yaml
jobs:
  global:
    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v4
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: eu-central-1

    - name: Setup Terraform
      uses: hashicorp/setup-terraform@v2

    - name: Setup Terrarium
      uses: terrarium-tf/github-action@vmaster

    - name: "default/foo stack"
      run: terrarium apply stage stacks/foo

```

### Bitbucket Pipeline

```yaml
image: hashicorp/terraform:1.5

definitions:
  caches:
    # Cache zip cli
    bins: /usr/bin

  steps:
    - step: &prepare
        name: install terrarium
        script:
          - wget https://github.com/terrarium-tf/cli/releases/download/v1.2.0/terrarium_1.2.0_linux_amd64.tar.gz
          - tar -xvzf terrarium_1.2.0_linux_amd64.tar.gz
          - mv ./terrarium /usr/bin/
        caches:
          - bins

    - step: &stack
        name: Apply Stack
        caches:
          - bins
        script:
          - export TF_IN_AUTOMATION=1
          - terrarium init $STAGE stacks/stack
          - terrarium apply $STAGE stacks/stack
        artifacts:
          - stacks/stack/*.tfplan

pipelines:
  branches:
    main:
      - step: *prepare
      - parallel:
          - step: *stack
```
## Cloud Providers

we support storing remote state for all 3 major Cloud Providers (AWS, GCP, Azure), if you dont want (or cant) use a remote state simply provide the
`--remote-state=false` option during the init command. 

> you can still configure your remote state by hand, but remember to deactive the automatic configuration (as above)

for providing sensitive variables use a local not checked in file `local.tfvars.json`

global variables should be stored in `global.tfvars.json`

stack specific variables should be stored in `app.tfvars.json`

### AWS

for [AWS](https://developer.hashicorp.com/terraform/language/settings/backends/s3) we configure the s3 bucket and the (optional) dynamo state locking from these variables:

* s3 bucket name : `tf-state-{project}-{region}-{account}`
* dynamo table : `terraform-lock-{project}-{region}-{account}`
* s3 file : `{name}.tfstate`
* the AWS credentials should be provided by your shell environment

### GCP

for [GCP](https://developer.hashicorp.com/terraform/language/settings/backends/gcs) we configure the bucket from these variables:

* bucket name : `tf-state-{project}`
* credentials: read from the `credentials` variable or from the environment variables `GOOGLE_BACKEND_CREDENTIALS` or `GOOGLE_CREDENTIALS`
* prefix: read from the `prefix` variable

### Azure

for [Azure](https://developer.hashicorp.com/terraform/language/settings/backends/azurerm) we configure the bucket from these variables:

* storage_account_name : read from the `account` variable
* resource_group_name : read from the `project` variable
* file (key) : `{name}.tfstate`
* container_name (bucket): `tf-state-{project}-{account}`

all optional variables from the documentation can be provided by their `ARM_*` environment variables (or through variables):

```go
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
```


## Development

Checkout the source and install golang dependencies with:

```shell script
$ go get
```

You should then be able to run a local version of the project with:

```shell script
$ go run main.go
```

To build and distribute the binary:

```shell script
$ goreleaser build --snapshot --clean
$ cp ./dist/terrarium_xxx/terrarium /usr/local/bin/terrarium
$ chmod a+x /usr/local/bin/terrarium
```
