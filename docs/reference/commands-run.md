---
layout: default
title: Run
parent: Commands
grand_parent: Reference
permalink: /reference/commands/run
---

{: .fs-10 .fw-300 }
# Run

{: .fs-6 .fw-300 }
The `run` command executes a request, routine or script at the path provided.

## Help Output

Run `nap run --help` to see information about the `run` command:

```
The run command executes a request, routine or script at the path provided.

Usage:
  nap run <target> [flags]

Flags:
  -e, --env path               add environment variables from a file path
  -h, --help                   help for run
  -p, --param <name>=<value>   add a single variable to the run as a <name>=<value> pair

Global Flags:
  -v, --verbose   verbose output
```

## Flags

### `--env` - Environment

Alias: `-e`. `string`. Optional.

Usage: `-e ./path/to/env.yml [-e ./path/to/another/env.yml] ...`

Add environment variables from a file path. To include multiple sets of variables, use the flag multiple times. Values found in files from later instances of the flag will overwrite those from earlier files.

Environments can be referenced by path or by just the file name. If a full path isn't provided, Nap will attempt to find the environment either in the directory of the target or near it. For example, given the above example, Nap will attempt to find `my-env` in the following locations:

* `./my-env.yml` - relative to the current working directory
* `./env/my-env.yml` - in an `env` folder, relative to the current working directory
* `./routines/my-env.yml` - in the target's directory
* `./routines/env/my-env.yml` - in an `env` folder within the target's directory

### `--param` - Parameter

Alias: `-p` `<name>=<value>`. Optional

Usage: `-p var1=val1 [-p var2=val2] ...`

Initialize a variable. To include multiple parameters, use the flag multiple times. If the same variable name is supplied multiple times, only the last value will be used. The `--param` flag will also overwrite values loaded via the `--env` flag.