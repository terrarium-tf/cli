# Terrarium

> a tiny wrapper for Terraform to make loading env vars transparent by convention

[![Coverage Status](https://coveralls.io/repos/github/terrarium-tf/cli/badge.svg?branch=main)](https://coveralls.io/github/terrarium-tf/cli?branch=main)

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


## Command

```
$ ./terrarium
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
terrarium [command]

Available Commands:
apply       A brief description of your command
completion  Generate the autocompletion script for the specified shell
destroy     A brief description of your command
help        Help about any command
import      A brief description of your command
init        A brief description of your command
plan        A brief description of your command
remove      A brief description of your command

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

checkout the source and install python dependencies with pip

```shell script
$ go get
```

```shell script
$ go run main.go
```

build the binary

```shell script
$ goreleaser build --snapshot
$ cp ./dist/cli_xxx/cli /usr/local/bin/terrarium
$ chmod a+x /usr/local/bin/terrarium
```
