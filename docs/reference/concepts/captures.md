---
layout: default
title: Captures
nav_order: 2
parent: Concepts
grand_parent: Reference
permalink: /reference/concepts/captures
---

{: .fs-10 .fw-300 }
# Captures

{: .fs-6 .fw-300 }
A capture stores a part of a response into a variable.

# Syntax

Each capture is made up of two parts: a `variable` and a `query`. These are separated by a colon (`:`) to form a key/value pair in the request YAML.

Here's an example of a valid captures section:

```yml
captures:
  myVar: jsonpath $.myVal
  myOtherVar: body
```

Let's break down the first capture into its parts:

```yaml
  myVar: jsonpath $.myVal
```

### Variable: `myVar`

The variable tells Nap where to store this capture. This variable will overwrite any previous value assigned to the same name, such as those supplied via an environment file or prior script or capture.

### Query: `jsonpath $.myVal`

The query tells Nap what part of the response we want to capture. This query will retrieve the `myVal` property from the root object in the repsonse body, assuming it is in JSON format.

{: .highlight }
For the full query reference, see [Concepts -> Queries](/reference/concepts/queries).