---
layout: default
title: Routines
nav_order: 1
parent: Reference
---

{: .fs-10 .fw-300 }
# Routines

{: .fs-6 .fw-300 }
The routine is Nap's most basic building block.

## Syntax

```yaml
kind: routine # required; defines the document as a routine
steps: # array; at least one step is required. 
  - run: ./request-1.yml # required; the path to the target to run
```

## Properties

### `kind` - Kind

`string`. Required. Allowed values: `routine`.

Defines the document type as a routine.

### `steps` - Steps to run

`array`. Required. Must contain at least one element.

Defines one or more requests or subroutines to execute.

### `steps[].run` - Path to run

`string`. Required. 

The path to the request or subroutine to execute.

## Requests

Requests are run in the order they appear in a routine. They are also run in the routine's main channel. In other words, requests will block further execution until completed, or in serial.

## Subroutines

A routine may run another routine. These are referred to as subroutines. A subroutine is run in its own channel. This allows multiple subroutines to be run in parallel.