---
layout: default
title: Predicates
nav_order: 6
parent: Reference
---

{: .fs-10 .fw-300 }
# Predicates

{: .fs-6 .fw-300 }
A predicate is a condition expression that evaluates a part of a response against an expectation. The result of a predicate is either `true` or `false`.

## Syntax

A predicate will always appear _between_ a query and a value (such as an assert's `expectation`). Here is a table of all supported predicates and their functions.

| Predicate    | Description                                 | Example                       |
|:-------------|:--------------------------------------------|:------------------------------|
| `==`         | Query and value are equal.                  | `status == 200`               |
| `!=`         | Query and value are not equal.              | `status != 200`               |
| `<`          | Query is less than value.                   | `duration < 2000`             |
| `<=`         | Query is less than or equal to value.       | `duration <= 2000`            |
| `>`          | Query is greater than value.                | `header Content-Length > 0`   |
| `>=`         | Query is greater than or equal to value.    | `jsonpath $.age >= 18`        |
| `contains`   | Query contains value as a `string`.         | `body contains Hello, World!` |
| `startswith` | Query begins with value as a `string`.      | `body startswith Hello`       |
| `endswith`   | Query ends with value as a `string`.        | `body endswith World!`        |
| `matches`    | Query matches value as a regular expression | `body endswith World!`        |
| `in`         | Query matches one of a set of values        | `status in [ 200, 201 ]`      |

## Negation

Any predicate can be negated to achieve the opposite effect by using the `not` keyword. For example:

```yml
status not in [ 400, 500 ] # will succeed as long as the status is NOT 400 or 500
```