---
layout: default
title: Scripts
nav_order: 10
parent: Reference
---

{: .fs-10 .fw-300 }
# Scripts

{: .fs-6 .fw-300 }
Nap supports ES5-compatible Javascript. The framework also provides some built-in objects and functions to facilitate integrating these scripts into requests and routines.

## Built-In Functions

### `nap.env.get()` - Get environment variable

Gets the value of an environment variable.

Syntax: 

```javascript
nap.env.get(key)
```

#### Parameters

* `key` - `string`. The name of the variable to get.

### `nap.env.set()` - Set environment variable

Sets the value of an environment variable.

Syntax: 

```javascript
nap.env.set(key, value)
```

#### Parameters

* `key` - `string`. The name of the variable to set.
* `value` - `string`. The value to assign.

### `nap.run()` - Run

Runs a provided file

Syntax: 

```javascript
nap.run(path)
```

#### Parameters

* `path` - `string`. The path to the target to run.

### `nap.fail()` - Fail

Trigger a failure; aborts the current routine.

Syntax: 

```javascript
nap.fail(message)
```

#### Parameters

* `message` - `string`. The message to display

## Built-In Data

### `nap.http` - HTTP data

`object`. Contains HTTP request and response data. Only set for pre- post- request scripts.

#### Properties

* `request` - `object`. HTTP request data. Contains the following properties:
  * `url` - `string`. The target URL.
  * `verb` - `string`. The request method.
  * `body` - `string`. The request body.
  * `headers` - `object`. The request headers. Contains properties and values that match the header names and values.
* `response` - `object`. HTTP response data. Contains `null` for pre-request scripts. Contains the following properties:
  * `statusCode` - `number`. The numeric HTTP status code (e.g. 200).
  * `status` - `string`. The string status code returned from the server.
  * `body` - `string`. The response body expressed as a string.
  * `jsonBody` - `object`. The response body expressed as an object. `null` if response Content-Type is not a JSON type.
  * `headers` - `object`. The response headers. Contains properties and values that match the header names and values.
  * `elapsedMs` - `number`. The duration of the request in milliseconds.

