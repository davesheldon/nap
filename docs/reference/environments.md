---
layout: default
title: Environments
nav_order: 8
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

{: .highlight }
For the `--env` command line reference, see [Reference -> Commands -> Run](/nap/reference/commands/run#--env---Environment).


{: .highlight }
For the full variable reference, see [Reference -> Variables](/nap/reference/variables).