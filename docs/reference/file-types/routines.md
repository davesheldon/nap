---
layout: default
title: Routines
nav_order: 2
parent: File Types
grand_parent: Reference
permalink: /reference/file-types/routines
---

{: .fs-10 .fw-300 }
# Routines

{: .fs-6 .fw-300 }
The routine is Nap's most basic building block.

## Syntax

```yml
kind: routine # required; defines the document as a routine
name: my routine # optional; used to identify this routine
env: # optional; variables to set before running this routine
  myvar: myval
steps: # array; at least one step is required. 
  - run: ./request-1.yml # required; the path to the target to run
    iterations: # optional; path(s) to variable iterations to run for this step.
```

## Properties

### `kind` - Kind

`string`. Required. Allowed values: `routine`.

Defines the document type as a routine.

### `name` - Name

`string` Optional.

A name used to identify the routine. This is used in any logs/output to refer to the routine. If a name isn't given, its file-name is used instead.

### `env` - Environment Variables

`object` Optional.

A set of variables to set before running the steps in this routine. Any number of variables may be included as YAML properties and values.

### `steps` - Steps to run

`array`. Required. Must contain at least one element.

Defines one or more requests or subroutines to execute.

### `steps[].run` - Step run path

`string`. Required. 

The path to the request or subroutine to execute.

### `steps[].iterations` - Step Iterations

`string | array`. Optional. 

One or more paths or globs pointing to environment files to load and iterate over for this step. If specified, the step will be run once for each environment file found. Each iteration is loaded on top of the existing set of variables. If no iterations are found then the step will run once with the normal environment.

## Requests

Requests are run in the order they appear in a routine. They are also run in the routine's main channel. In other words, requests will block further execution until completed, or in serial.

## Subroutines

A routine may run another routine. These are referred to as subroutines. A subroutine is run in its own channel. This allows multiple subroutines to be run in parallel.