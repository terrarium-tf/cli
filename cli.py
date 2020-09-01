import json
import logging
import os
import sys
import click

from python_terraform import Terraform, IsFlagged


def config_file(file='global.tfvars.json') -> dict:
    """read a config file"""

    if os.path.exists(file):
        click.echo("using config {}".format(os.path.abspath(file).replace(os.path.abspath(os.path.curdir)+"/", '')))
        with open(file) as json_file:
            data = json.load(json_file)

            if type(data) is dict:
                return data

            raise Exception('not a valid file')

    return {}


def initialize_tf(stack) -> Terraform:
    """builds terraform wrapper"""
    # append the global vars as TF_VAR env vars (so we dont raise a warning)
    for k, value in config_file().items():
        if type(value) is list or type(value) is dict:
            value = json.dumps(value)
        elif type(value) is bool:
            value = "1" if value is True else "0"
        elif type(value) is str:
            pass
        else:
            value = str(value)
        os.environ['TF_VAR_{}'.format(k)] = value

    return Terraform(working_dir=stack)


def execute(stack: str, command: str, args: list, var:dict = {}):
    t = initialize_tf(stack)

    return_code, stdout, stderr = t.cmd(command, *args, capture_output=False, no_color=None, **var)

    sys.exit(return_code)


def workspaced_command(stack: str, environment: str, command: str, args: list, var_file = None):
    logging.basicConfig(filename='example.log', level=logging.DEBUG)

    t = initialize_tf(stack)

    vars = {
        **config_file("{}/app.tfvars.json".format(stack)),
        **config_file("{}/default.tfvars.json".format(stack)),
        'environment': environment
    }

    t.cmd('workspace', 'new', environment, no_color=None)
    t.cmd('workspace', 'select', environment, no_color=None, capture_output=False)

    return_code, stdout, stderr = t.cmd(command, *args, capture_output=False, no_color=None, var=vars)

    sys.exit(return_code)


@click.command()
@click.option('--initial', default=False, help="really initial with state infrastructure creation", is_flag=True)
@click.argument('stack', required=True)
def init(initial, stack):
    """initialize a given stack"""
    global_vars = config_file()
    app_vars = {}
    options = {}

    if os.path.exists("{}/app.tfvars.json".format(stack)):
        app_vars = {
            **app_vars,
            **config_file("{}/app.tfvars.json".format(stack))
        }

    if os.getenv('CI') is None:
        options = {
            **options,
            'verify_plugins': True,
            'upgrade': True,
            'force_copy': IsFlagged,
            'reconfigure': IsFlagged,
        }

    if initial is False:
        options['backend_config'] = [
            "region={}".format(global_vars['region']),
            "dynamodb_table={}-terraform-lock".format(global_vars['project']),
            "bucket={}-tf-{}-{}".format(global_vars['project'], global_vars['account_id'], global_vars['region']),
            "key={}.tfstate".format(app_vars['name'])
        ]

    execute(stack, 'init', ['.'], options)


@click.command()
@click.argument('environment', required=True)
@click.argument('stack', required=True)
def apply(environment, stack):
    """apply a given stack"""
    t = initialize_tf(stack)

    vars = {
        **config_file("{}/app.tfvars.json".format(stack)),
        **config_file("{}/default.tfvars.json".format(stack)),
        'environment': environment
    }

    t.cmd('workspace', 'new', environment, no_color=None)
    t.cmd('workspace', 'select', environment, no_color=None, capture_output=False)

    if os.getenv('CI') is not None:
        plan_file = "{}-{}.plan".format(stack.replace('/', '-'), environment)
        t.plan(capture_output=False, var=vars, no_color=None, out=plan_file)
        return_code, stdout, stderr = t.apply(plan_file, capture_output=False, no_color=None, var=None, var_file=None)
    else:
        return_code, stdout, stderr = t.apply(capture_output=False, var=vars, no_color=None)

    sys.exit(return_code)


@click.command()
@click.argument('environment', required=True)
@click.argument('stack', required=True)
def destroy(environment, stack):
    """destroy a given stack"""
    workspaced_command(stack, environment, 'destroy', [])


@click.command()
@click.argument('stack', required=True)
def fmt(stack):
    """format a given stack"""
    execute(stack, 'fmt', [])


@click.command(name="import")
@click.argument('environment', required=True)
@click.argument('stack', required=True)
@click.argument('resource', required=True)
@click.argument('id', required=True)
def import_cmd(environment, stack, resource, id):
    """import a given stack"""
    workspaced_command(stack, environment, 'import', [resource, id])


@click.command()
@click.argument('environment', required=True)
@click.argument('stack', required=True)
@click.argument('resource', required=True)
def rm(environment, stack, resource):
    """import a given stack"""
    workspaced_command(stack, environment, 'state', ['rm'], resource ,var_file=None)

@click.command()
@click.argument('environment', required=True)
@click.argument('stack', required=True)
@click.argument('resource', required=True)
def taint(environment, stack, resource):
    """import a given stack"""
    workspaced_command(stack, environment, 'taint', [resource])
