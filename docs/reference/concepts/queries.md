---
layout: default
title: Queries
nav_order: 3
parent: Concepts
grand_parent: Reference
permalink: /reference/concepts/queries
---

{: .fs-10 .fw-300 }
# Queries

{: .fs-6 .fw-300 }
A query retrieves a part of an HTTP response.

## Syntax

Each query begins with a keyword. Here is a table of all supported keywords and their functions.

| Keyword     | Description                                       |
|:------------|:--------------------------------------------------|
| `body`      | The raw HTTP response body, expressed as a string |
| `duration`  | The HTTP execution duration in milliseconds       |
| `header`    | The value of an HTTP response header              |
| `cookie`    | The value of an HTTP response cookie              |
| `jsonpath`  | The result of a jsonpath expression               |
| `status`    | The numeric HTTP response status code             |

## Filters

Some queries (such as `header` and `jsonpath`) also require a filter. The filter appears after the query keyword and is separated from the keyword by a space. For example, to query the `Content-Type` header, use the query: `header Content-Type`.

## Query Examples

Below, each query is explained in greater detail. Variable captures are provided as examples.

### `body` - HTTP Response Body

Supported Format(s): `body`

The entire body from the HTTP response in text format.

```yaml
captures:
  bodyText: body
```

### `duration` - HTTP Execution Duration

Supported Format(s): `duration`

The duration (in milliseconds) that the HTTP response took to receive.

```yaml
captures:
  elapsedMs: duration
```

### `header` - HTTP Response Header

Supported Format(s): `header name`

A single header's value from the HTTP response. This query requires a header name to be supplied as a filter.

```yaml
captures:
  contentType: header Content-Type
```

### `cookie` - HTTP Response Cookie

Supported Format(s): `cookie name`, `cookie name[attribute]`

A single cookie's value or attribute from the HTTP response. This query requires a cookie name to be supplied as a filter. It also accepts an optional attribute. This attribute appears at the end of the cookie name, surrounded by square brackets (`[]`).

```yaml
captures:
  myCookieVal: cookie my_cookie

  # full list of supported attributes
  myCookieVal2: cookie my_cookie[Value] # same as the simpler form "cookie my_cookie"
  myCookieExpires: cookie my_cookie[Expires]
  myCookieMaxAge: cookie my_cookie[Max-Age]
  myCookieDomain: cookie my_cookie[Domain]
  myCookiePath: cookie my_cookie[Path]
  myCookieSecure: cookie my_cookie[Secure]
  myCookieHttpOnly: cookie my_cookie[HttpOnly]
  myCookieSameSite: cookie my_cookie[SameSite]
```

### `jsonpath` - JSONPath Expression

Supported Format(s): `jsonpath <expression>`

The result of a [JSONPath](https://ietf-wg-jsonpath.github.io/draft-ietf-jsonpath-base/draft-ietf-jsonpath-base.html) expression queried against a JSON response body.

```yaml
captures:
  firstId: $[0].id
  resultName: $.name
  arrayLength: $[*].length()
  filteredId: $[?(@.name == "John Smith")].id
```

### `status` - HTTP Response Status

Supported Format(s): `status`

The numeric status code received from the HTTP response

```yaml
captures:
  statusCode: status
```
