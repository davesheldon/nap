---
layout: default
title: Variables
nav_order: 5
parent: Concepts
grand_parent: Reference
permalink: /reference/concepts/variables
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
For the full environment reference, see [File Types -> Environments](/reference/file-types/environments).

## Inline

Individual variables may be added via the `run` command's `--param` flag.

{: .highlight }
For the `--param` command line reference, see [Commands -> Run](/reference/commands/run#--param---parameter).

## In Scripts

Variables may be read and mutated from scripts via `nap.env.get(key)` and `nap.env.set(key, value)`.

{: .highlight }
For the full script reference, see [File Types -> Scripts](/reference/file-types/scripts).