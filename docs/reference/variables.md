---
layout: default
title: Variables
nav_order: 9
parent: Reference
---

{: .fs-10 .fw-300 }
# Variables

{: .fs-6 .fw-300 }
Variables store information that can be used in routines, requests and scripts.

## Syntax

```yml
# request-1.yml
# statusAssert: status in [ 200, 201 ]
kind: request
name: Cat Breeds - Assertion/Capture Testing
path: https://catfact.ninja/breeds
asserts:
  - ${statusAssert} 
```

Variables are injected into requests and routines whenever they're loaded. Variables are reference by name, in the format `${variable}`.

## In Environments

Variables may be added en masse via environment files.

{: .highlight }
For the full environment reference, see [Reference -> Environments](/reference/environments).

## Inline

Individual variables may be added via the `run` command's `--param` flag.

{: .highlight }
For the `--param` command line reference, see [Reference -> Commands -> Run](/reference/commands/run#--param---Parameter).

## In Scripts

Variables may be read and mutated from scripts via `nap.env.get(key)` and `nap.env.set(key, value)`.

{: .highlight }
For the full script reference, see [Reference -> Scripts](/reference/scripts).