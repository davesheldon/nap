![nap logo](https://user-images.githubusercontent.com/7782805/147585754-36a405c8-5821-4482-963e-8a8eb75cd8e6.png)

# Nap

Nap is a _FAST_, file-based command-line interface (CLI) for executing HTTP tests, using plain old YAML. 

Here's an example of a request in Nap:

`./request-1.yml`
```yml
kind: request
name: Cat Breeds - Assertion/Capture Testing
path: https://catfact.ninja/breeds
asserts: # failed asserts go to stderr
  - status == 200
  - duration < 1000
  - header Content-Type == application/json
  - jsonpath $.data[0].breed matches ^Abyss.+
captures: # captures can be used in later requests using the ${variable} syntax
  firstBreed: jsonpath $.data[0].breed
  secondBreed: jsonpath $.data[1].breed
```

Let's run it:

```bash
$ nap run ./request-1.yml
Run finished in 297ms. 1/1 succeeded.
```

Easy!

What makes Nap _FAST_? Routines make Nap _FAST_. Here's what a routine looks like in Nap:

`./routine.yml`
```yml
kind: routine
name: main routine
steps:
  - run: request-1.yml
  - run: request-2.yml
  - run: sub-routine-1.yml # yep, you can call routines from routines.
  - run: sub-routine-2.yml # additional sub-routines run in PARALLEL 😱
```

We run routines the same way we run requests:

```bash
$ nap run ./routine.yml
Run finished in 699ms. 4/4 succeeded.
```

# Installation Options

Ready to install? Here are your options.

## Windows .exe

64-bit zip archive: [Download](https://raw.githubusercontent.com/davesheldon/nap/main/dist/windows/amd64/nap.zip)

## Using go get

```bash
$ go install github.com/davesheldon/nap@latest
```

## Building the Source

```bash
$ git clone https://github.com/davesheldon/nap.git
$ cd nap
$ go install
```


# Getting Started

Follow these steps to get to work (or just start creating files and running them, Nap won't be mad).

## Starting a New Project

To create a new project (not required, but enables some QoL features, such as templates), run the `new` command:

```bash
$ nap new my-project
Created a new project called my-project. Run cd my-project to get started.
```

## Project Structure

A new project consists of the project directory and a few default folders and files:

```
/my-project/
  /.template/
    request.yml
    routine.yml
    env.yml
    script.js
  /env/
    default.yml
  /requests/
    request-1.yml
  /routines/
    routine-1.yml
  /scripts/
    script-1.js
```

## Components

A Nap project may consist of several different components: requests, environments, scripts and routines.

# Requests

A **request** represents a single HTTP request. To generate a request, use the `generate` command:

```bash
$ nap generate request requests/my-request.yml
```

By default, this creates a file called `my-request.yml` inside the `requests` folder like the following:

```yml
kind: request
name: Cat Facts
path: https://catfact.ninja//facts
verb: GET
body:
headers:
  - Accept: application/json
```

**Note:** to customize the default request template, edit `request.yml` found inside the `.templates` folder.

## Running Requests

To run a specific request, use the `run` command as follows:

```bash
$ nap run requests/my-request.yml
- requests/my-request.yml: 200 OK
```

## Asserts

Nap supports a variety of assertions. Here is an example containing all the possible usages:

```yml
kind: request
name: Cat Breeds - Assertion Example
path: https://catfact.ninja/breeds
asserts:
  - status == 200
  - duration > 0
  - duration < 1000
  - duration >= 0
  - duration <= 1000
  - header Content-Type == application/json
  - jsonpath $.data[0].breed matches ^Abyss.+
  - jsonpath $.data[0].breed contains byssi
  - jsonpath $.data[0].breed startswith Ab
  - jsonpath $.data[0].breed endswith ian
```

Asserts are broken up into 3 parts: query, predicate and expectation. Consider the example:

```yml
jsonpath $.data[0].breed matches ^Abyss.+
```

In this example, `jsonpath $.data[0].breed` is the query, `matches` is the predicate and `^Abyss.+` is the expectation.

## Captures

Captures allow us to evaluate a `query` against the response and save it to a variable. Here is an example of the usage:

```yml
kind: request
name: Cat Breeds - Capture Example
path: https://catfact.ninja/breeds
captures:
  myStatus: status
  myDuration: duration
  firstBreed: jsonpath $.data[0].breed
  secondBreed: jsonpath $.data[1].breed
```

This is especially useful for test runs where an auth token needs to be requested and re-used or where an entity is created and an Id needs to be stored.

# Environments

The `env` folder contains a default YAML configuration file: `default.yml`. By default, this file is empty. Values added to the default configuratoin may be substituted within requests or routines. Here is an example of our first request with the base URL stored as a variable:

requests/my-request.yml:
```yml
kind: request
name: Cat Facts
path: ${baseurl}/facts
verb: GET
body: ""
headers:
    - Accept: application/json
```

env/default.yml:
```yml
baseurl: https://catfact.ninja
```

You may create new configurations either by adding a .yml file manually to the `env` folder or via the `generate` command:

```bash
$ nap generate env env/my-env.yml
```

**Note:** to customize the default environment template, edit `env.yml` found inside the `.templates` folder.

To run a request with a particular set of environment variables, use the `run` command with the `--env` or `-e` flag:

```bash
$ nap run requests/my-request.yml -e env/my-env.yml
- requests/my-request.yml: 200 OK
```

# Scripts

A **script** is a file containing [ES5-compatible](https://www.w3schools.com/js/js_es5.asp) Javascript. Nap supports ES6 Javascript via the [Otto](https://github.com/robertkrimen/otto) library, which means the same [limitations](https://github.com/robertkrimen/otto) are in play as mentioned on the Otto project page.

Scripts can be run before or after a request:

request.yml:
```yml
kind: request
name: Cat Facts
path: ${baseurl}/facts
verb: GET
pre-request-script-file: ../scripts/script-1.js
post-request-script-file: ../scripts/script-2.js
headers:
    - Accept: application/json
```

Scripts may also be inlined:

request.yml:
```yml
kind: request
name: Cat Facts
path: ${baseurl}/facts
verb: GET
pre-request-script: |
  console.log('Hello, World!');
post-request-script: |
  console.log('Goodbye, World!');
headers:
    - Accept: application/json
```

## Built-in functions

Nap provides several built-in functions for scripts to use. These are all nested under the global variable `nap`:

| Function | Description |
|-|-|
| nap.env.get(key: string) | Returns the value of an environment variable |
| nap.env.set(key: string, value: string) | Sets the (in-memory) value of an environment variable |
| nap.run(path) | Locates the referenced file, resolves its type and runs it |
| nap.fail(message: string) | Trigger a failure with a message; abort the rest of the routine |

## Environment Variables in Scripts

The templating syntax supported for environments is not supported in scripts. In order to access environment variables from inside a script, you must use the built-in functions.

# Routines

A **routine** is a file containing instructions for running one or more requests, scripts and/or subroutines in a specific order. The routine is _the_ first-class unit of execution in Nap. In fact, even requests and scripts that are run directly are first inserted into a routine at runtime in order to be executed. To generate a routine, use the `generate` command:

```bash
$ nap generate routine routines/my-routine.yml
```

By default, this creates a file called `my-routine.yml` inside the `routines` folder like the following:

```yml
kind: routine
name: Basic Routine
steps:
  - run: ../requests/my-request.yml
```

Each step may specify a target to run. Paths are relative to routine file location, so in this case to we must back out of the routines folder and into the requests folder to run out request.

**Note:** to customize the default routine template, edit `routine.yml` found inside the `.templates` folder.

## Running Routines

To run a routine, use the `run` command as follows:

```bash
$ nap run routines/my-routine.yml
- ../requests/my-request.yml: 200 OK
```

## Subroutines

A **subroutine** is a routine step that runs another routine. In this way you may use a single routine to run entire suites of tests:

```yml
kind: routine
name: Routine With Subroutines
steps:
  - run: my-routine.yml
  - run: my-other-routine.yml
```

Each subroutine will run within its own goroutine. This allows designing each subroutine as an end-to-end integration test that can run in parallel to other tests. 


# Concurrency

Nap is built upon a concurrency model where, by default, each routine runs in a separate thread. 

## Environment Variables

Each routine-thread in Nap receives a snapshot of the latest set of environment variables (including any changes made via scripts on the parent routine up until that point). This allows for scenarios such as performing authentication up-front, setting a token as a variable, and then running multiple routines in parallel that rely on that token. For example:

```yml
kind: routine
name: Authenticated Routine
steps:
  - run: ../requests/auth.yml
  - run: authenticated-routine-1.yml
  - run: authenticated-routine-2.yml
  - run: authenticated-routine-3.yml
```

The above results in a workflow like the following:

1. Run the auth request (sets the auth token into env)
2. Start the remaining subroutines (each receives its own copy of the current context, including a fresh scripting instance and snapshot of the environment variables)
3. Each routine runs, using its copy of the env auth token.

# Scripting Considerations

Since Nap can run scripts directly, entire workflows can be orchestrated using them. This format is encouraged for more advanced workflows. For example:

```javascript
var start = function(){
  console.log("Starting script"); // logs will show in the terminal window

  nap.run("../requests/auth.yml");

  if (nap.env.get("auth_token") && nap.env.get("auth_token").length > 0){
    console.log("Authenticated.");
  }
  else {
    nap.fail("Authentication failed.");
    return;
  }

  runRoutines();
};

var runRoutines = function(){
  console.log("Running routines synchronously");

  // Since each routine is being started via a script, they'll run in serial.
  // Subroutines of these routines will still run in parallel.

  nap.run("../routines/routine-1.yml"); 

  var errorMessage = nap.env.get('error_mesage');

  if (errorMessage && errorMessage.length > 0) {
    nap.fail(errorMessage);
    return;
  }

  nap.run("../routines/routine-2.yml"); 
  nap.run("../routines/routine-3.yml");
};

start();
```