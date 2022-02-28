# Terrarium

> a tiny wrapper for Terraform to make loading env vars transparent by convention

[![Coverage Status](https://coveralls.io/repos/github/terrarium-tf/cli/badge.svg?branch=main)](https://coveralls.io/github/terrarium-tf/cli?branch=main)

**Builds Terraform Commands, easing those steps:**

* collects defined var-files
* switches to the given workspace (can create new one)
* runs the given terraform command with the multiple -var-files options in correct order.

## Setup

download the binary matching your OS from [here](https://github.com/terrarium-tf/cli/releases)

## Anatomy of Stacks and Configs

```
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
$ ./terrarium
Builds Terraform Commands, easing those steps:
* collects defined var-files
* switches to the given workspace (can create new one)
* runs the given terraform command with the multiple -var-files options in correct order.

You can override the default terraform binary with "-t"
Add "-v" for more verbose logging.

Usage:
  terrarium [command]

Examples:
terrarium [command] workspace path/to/stack -v -t echo

Available Commands:
  apply       Apply a given Terraform Stack
  completion  Generate the autocompletion script for the specified shell
  destroy     Destroy a given Terraform stack
  help        Help about any command
  import      Import a remote resource into a local terraform resource
  init        initializes a stack with optional remote state
  plan        Creates a diff between remote and local state and prints the upcoming changes
  remove      Removes a remote resource from the terraform state

Flags:
  -h, --help               help for terrarium
  -t, --terraform string   terraform binary found in your path (default "/usr/local/bin/terraform")
  -v, --verbose            display extended informations

Use "terrarium [command] --help" for more information about a command.
```


## Use within Github-Actions

```yaml
jobs:
  global:
    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v1
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: eu-central-1

    - name: Setup Terraform
      uses: hashicorp/setup-terraform@v1

    - name: Setup Terrarium
      uses: terrarium-tf/github-action@vmaster

    - name: "default/foo stack"
      run: terrarium apply stage stacks/foo

```

## Development

checkout the source and install golang dependencies with pip

```shell script
$ go get
```

```shell script
$ go run main.go
```

build the binary

```shell script
$ goreleaser build --snapshot --rm-dist
$ cp ./dist/terrarium_xxx/terrarium /usr/local/bin/terrarium
$ chmod a+x /usr/local/bin/terrarium
```
