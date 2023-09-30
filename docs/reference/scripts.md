---
layout: default
title: Scripts
nav_order: 7
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
