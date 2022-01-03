![nap logo](https://user-images.githubusercontent.com/7782805/147585754-36a405c8-5821-4482-963e-8a8eb75cd8e6.png)

# Nap

Nap is a _FAST_, file-based framework for creating and running integration tests over HTTP.

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

# Components

A Nap project may consist of several different components: requests, environments, scripts and routines.

## Requests

A **request** represents a single HTTP request. To generate a request, use the `generate` command:

```bash
$ nap generate request my-request.yml
```

By default, this creates a YAML file inside the `requests` folder like the following:

```yml
type: request
path: https://cat-fact.herokuapp.com/facts
verb: GET
body:
headers:
  - Accept: application/json
```

**Note:** to customize the default request template, edit `request.yml` found inside the `.templates` folder.

### Running Requests

To run a specific request, use the `run` command as follows:

```bash
$ nap run requests/my-request.yml
- requests/my-request.yml: 200 OK
```

## Environments

The `env` folder contains a default YAML configuration file: `default.yml`. By default, this file is empty. Values added to the default configuratoin may be substituted within requests or routines. Here is an example of our first request with the base URL stored as a variable:

requests/my-request.yml:
```yml
type: request
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

## Scripts

TODO

## Routines

A **routine** is a file containing instructions for running one or more requests, scripts or assertions in a specific order. To generate a routine, use the `generate` command:

```bash
$ nap generate routine my-routine
```

By default, this creates a YAML file inside the `routines` folder like the following:

```yml
type: routine
steps:
  - run: ../requests/my-request.yml
```

Each step may specify a target to run. Paths are relative to routine file location, so in this case to we must back out of the routines folder and into the requests folder to run out request.

**Note:** to customize the default routine template, edit `routine.yml` found inside the `.templates` folder.

### Running Routines

To run a routine, use the `run` command as follows:

```bash
$ nap run routines/my-routine.yml
- ../requests/my-request.yml: 200 OK
```

### Subroutines

A **subroutine** is a routine step that runs another routine. In this way you may use a single routine to run entire suites of tests:

```yml
type: routine
steps:
  - run: my-routine.yml
  - run: my-other-routine.yml
```

Each subroutine will run within its own goroutine. This allows designing each subroutine as an end-to-end integration test that can run in parallel to other tests. 
