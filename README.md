# Nap

Nap is a file-based framework for automating the execution of config-driven HTTP requests and scripts.

# Installation Options

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

Follow these steps to get to work.

## Starting a New Project

To create a new project, run the `new` command:

```bash
$ nap new my-project
Created a new project called my-project. Run cd my-project to get started.
```

# Requests

A **request** represents a single HTTP request. To generate a request, use the `generate` command:

```bash
$ nap generate request my-request.yml
```

By default, this creates a YAML file inside the `requests` folder like the following:

```yml
name: my-request
path: https://cat-fact.herokuapp.com/facts
verb: GET
body:
headers:
    - Accept: application/json
```

**Note:** to customize the default request template, edit `request.yml` found inside the `.templates` folder.

## Running Requests

To run a specific request, use the `run` command as follows:

```bash
$ nap run request my-request
- my-request.yml: 200 OK
```

# Routines

A **routine** is a file containing instructions for running one or more requests, scripts or assertions in a specific order. To generate a routine, use the `generate` command:

```bash
$ nap generate routine my-routine
```

By default, this creates a YAML file inside the `routines` folder like the following:

```yml
name: my-routine
run:
    - type: request 
      name: my-request
      expectStatusCode: 200
```

**Note:** to customize the default routine template, edit `routine.yml` found inside the `.templates` folder.

## Expectations

An **expectation** defines an attribute to be tested after the request is executed. Nap supports several ways to help determine whether a request meets expectations or not:

### Status Code Expectation

```yml
expectStatusCode: 2XX
```

The **status code** expectation will test against the expected status code. An `X` will match any digit in the response. For example, to match against all successful response codes, the expectation would be `2XX`.

### Headers Expectation

```yml
expectHeaders:
  - Content-Type: application/json
  - X-CUSTOM-HEADER: custom-value
```

The **headers** expectation tests to ensure the provided headers exist with the expected values. It isn't a strict comparison against _all_ headers.

### Response Content Expectation

```yml
expectResponseContent: |
  {"status":"success"}
```

The **response content** expectation performs a strict match against the string-encoded content of the response.

### JSON Expectation

```yml
expectJson:
  status: success
```

The **JSON** expectation will test specific parts of a JSON-encoded response object. The above example would match against an object literal such as:

```json
{
    "status": "success",
    "result": [
        {
            "name": "example"
        }
    ]
}
```

## Subroutines

To work with multiple routines, use the subroutine pattern:

```yml
run:
    - type: routine
      name: my-routine
    - type: routine
      name: my-other-routine
```

## Running Routines

To run a routine, use the `run` command as follows:

```bash
$ nap run routine my-routine
- my-request.yml: 200 OK
```

# Environment Variables

The `env` folder contains a default YAML configuration file: `default.yml`. By default, this file is empty. Values added to the default configuratoin may be substituted within requests or routines. Here is an example of our first request with the base URL stored as a variable:

requests/my-request.yml:
```yml
name: my-request
path: ${baseurl}/facts
verb: GET
type: request
body: ""
headers:
    - Accept: application/json
```

env/default.yml:
```yml
baseurl: https://cat-fact.herokuapp.com
```

You may create new configurations either by adding a .yml file manually to the `env` folder or via the `generate` command:

```bash
$ nap generate env my-env
```

To run a request with a particular environment, use the `run` command with the `--env` or `-e` flag:

```bash
$ nap run request my-request -e my-env
- my-request.yml: 200 OK
```

