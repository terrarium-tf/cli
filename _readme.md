# Terrarium

> a tiny wrapper for Terraform to make loading env vars transparent by convention

## Setup

download the binary matching your OS from [here](https://github.com/terrarium-tf/cli/releases)

## Anatomy of Stacks and Configs

```
 |- global.tfvars.json # variables available to all stacks (relative to cwd)
 |
 | - stacks
    |
    | - foo
        | - app.tfvars.json # default variables for this stack
        | - default.tfvars.json # variables for the default workspace (environment)
        | - stage.tfvars.json # variables for the stage workspace (environment)
        |
        | - main.tf your stack entrypoint
```

## Command

```
$ ./terrarium
Usage: terrarium [OPTIONS] COMMAND [ARGS]...

  terrarium cli

Options:
  --version  Show the version and exit.
  --help     Show this message and exit.

Commands:
  apply    apply a given stack
  destroy  destroy a given stack
  fmt      format a given stack
  import   import a given stack
  init     initialize a given stack
  rm       import a given stack
  taint    import a given stack
```

### apply

rollout a stack

```shell script
$ ./terraform apply default stacks/foo
using config global.tfvars.json
using config stacks/foo/app.tfvars.json
using config stacks/foo/default.tfvars.json
```

another environment (terraform workspace)
```shell script
$ ./terraform apply stage stacks/foo
using config global.tfvars.json
using config stacks/foo/app.tfvars.json
using config stacks/foo/stage.tfvars.json
```

### destroy

destroys a stack

```shell script
$ ./terraform destroy default stacks/foo
```

### fmt

terraform format

```shell script
$ ./terraform fmt stacks/foo
```

### import

import an existing resource into the stack

```shell script
$ ./terraform import default stacks/foo aws_ecs_cluster.cluster ARN
```

### rm

removes an existing resource from the state

```shell script
$ ./terraform rm default stacks/foo aws_ecs_cluster.cluster
```


### taint

taints a specific resource for recreation

```shell script
$ ./terraform taint default stacks/foo aws_ecs_cluster.cluster
```

### init

initializes a stack

```shell script
$ ./terraform init stacks/foo
```

initializes a stack (without a remote backend, probably only once at the very first beginning)

```shell script
$ ./terraform init --initial stacks/foo
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
      run: terrarium apply default stacks/foo

```

## Development

checkout the source and install python dependencies with pip

```shell script
$ pip3 -r requirements.txt
```

use the python script

```shell script
$ ./terrarium
```

build the binary

```shell script
$ python3 -m nuitka --follow-imports ./terrarium
$ cp ./terrarium.bin /usr/local/bin/terrarium
$ chmod a+x /usr/local/bin/terrarium
```

## TODO

* tests
