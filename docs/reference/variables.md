---
layout: default
title: Variables
nav_order: 8
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

## In Scripts

Variables may also be referenced in scripts via `nap.env.get(variable)`.

{: .highlight }
For the full script reference, see [Reference -> Scripts](/nap/reference/scripts).