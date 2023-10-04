---
layout: default
title: Asserts
nav_order: 4
parent: Reference
---

{: .fs-10 .fw-300 }
# Asserts

{: .fs-6 .fw-300 }
Asserts are how we evaluate expectations against an HTTP response.

## Components

Each Assert is made up of three main parts: a `query`, a `predicate` and an `expectation`. These parts combine to allow you to write a wide range of powerful tests.

Here's an example of an assert that ensures that HTTP response status is equal to 200:

```yml
status == 200
```

Let's break this down into the three parts.

### Query: `status`. 

The query tells Nap what part of the response we want to validate. In this case we're checking the HTTP response status.

{: .highlight }
For the full query reference, see [Reference -> Queries](/reference/queries).

### Predicate: `==`

The predicate tells Nap what sort of validation we want this assert to perform. The `==` means we're testing for equality.

{: .highlight }
For the full predicate reference, see [Reference -> Predicates](/reference/predicates).

### Expectation: `200`

The expectation is the specific value we're testing against. For this assert to succeed, the result of our query (`status`) must equal the expectation (`200`).

## Failures

Assert failures are written to `stderr` and follow the form:
```
[ERROR] <request name>: Assert failed "<query> => <actual> <predicate> <expectation>"
```

For example:
```
[ERROR] Cat Breeds - Assertion/Capture Testing: Assert failed "jsonpath $.data[0].breed => Abyssinian matches ^Ayss.+"
```