---
layout: default
title: The Basics
nav_order: 2
---

{: .fs-10 .fw-300 }
# The Basics

{: .fs-6 .fw-300 }
Nap has a few basic building-blocks that you can leverage to quickly write your first test.

## Requests

Here's an example of a request in Nap:

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

Let's save this to a file and run it. Let's call it `request-1.yml`. To run this request, we'll go to the command-line:

```bash
$ nap run ./request-1.yml
Run finished in 297ms. 1/1 succeeded.
```

Easy enough. For the full spec, containing all the options, see [Reference -> Request](/nap/reference/requests)

For now, let's move on to routines.

## Routines

What makes Nap _FAST_? Routines make Nap _FAST_. Here's what a routine looks like in Nap:

```yml
kind: routine
name: main routine
steps:
  - run: request-1.yml
```

We can save and run a routine the same way we would run a request. Let's call this one `routine-1.yml` and run it:

```bash
$ nap run ./routine-1.yml
Run finished in 312ms. 1/1 succeeded.
```

Oops. That was actually slower. Let's change this up a bit. Let's start by adding a second request, called `request-2.yml`:

```yml
kind: request
name: Cat Facts - Succeeds
path: https://catfact.ninja/facts
asserts:
  - status == 200
  - duration < 1000
```

Simple enough. Now let's add a couple more routines. These will be redundant, but we're just illustrating a point here. We'll call them `subroutine-1.yml` and `subroutine-2.yml`, respectively:

```yml
kind: routine
name: example subroutine 1
steps:
  - run: request-1.yml
  - run: request-2.yml
```

```yml
kind: routine
name: example subroutine 2
steps:
  - run: request-2.yml
  - run: request-1.yml
```

Finally, we'll update `routine-1.yml` to the following:

```yml
kind: routine
name: main routine
steps:
  - run: subroutine-1.yml # yep, you can call routines from routines.
  - run: subroutine-2.yml # additional sub-routines run in PARALLEL 😱
```

And run it:

```bash
$ nap run ./routine-1.yml
Run finished in 326ms. 4/4 succeeded.
```

This is the power of Nap. Since `subroutine-1.yml` and `subroutine-2.yml` aren't blocking each other, we can write a very large set of tests and have many of them running in parallel, making the whole test run much, much faster than solutions such as Postman or a regular bash script using `curl`.

{: .warning }
Be sure not to reference a routine from itself. Unless you like stalling forever, of course.

You can view the full code for example on Github [here](https://github.com/davesheldon/nap/tree/main/examples/routines/basic).