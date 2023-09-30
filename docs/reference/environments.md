---
layout: default
title: Environments
nav_order: 7
parent: Reference
---

{: .fs-10 .fw-300 }
# Environments

{: .fs-6 .fw-300 }
Environments allow us to initialize one or more variables in Nap at the beginning of a run.

## Syntax

```yml
myVar: myVal
myOtherVar: myOtherVal
```

An environment file is a `.yml` file arranged into key/value pairs. 

During Nap's initialization, each key will be saved to a variable with its corresponding value.

## Usage

Use the `-e` or `--env` flag to include an environment file.

### Example

```bash
$ nap run ./routines/routine-1.yml -e my-env
```

### Remarks

Environments can be referenced by path or by just the file name. If a full path isn't provided, Nap will attempt to find the environment either in the directory of the target or near it. For example, given the above example, Nap will attempt to find `my-env` in the following locations:

* `./my-env.yml` - relative to the current working directory
* `./env/my-env.yml` - in an `env` folder, relative to the current working directory
* `./routines/my-env.yml` - in the target's directory
* `./routines/env/my-env.yml` - in an `env` folder within the target's directory