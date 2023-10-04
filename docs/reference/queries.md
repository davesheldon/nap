---
layout: default
title: Queries
nav_order: 6
parent: Reference
---

{: .fs-10 .fw-300 }
# Queries

{: .fs-6 .fw-300 }
A query describes a part of an HTTP response.

## Syntax

Each query begins with a keyword. Here is a table of all supported keywords and their functions.

| Keyword | Description |
|:--------|:------------|
| `body` | The raw HTTP response body, expressed as a string |
| `duration` | The HTTP execution duration in milliseconds |
| `header` | The value of an HTTP response header |
| `jsonpath` | The result of a jsonpath expression |
| `status` | The numeric HTTP response status code |

## Filters

Some queries (such as `header` and `jsonpath`) also require a filter. The filter appears after the query keyword and is separated from the keyword by a space. For example, to query the `Content-Type` header, use the query: `header Content-Type`.