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
# request-1.yml
kind: request
name: Cat Breeds - Assertion/Capture Testing
path: https://catfact.ninja/breeds
asserts: # failed asserts go to stderr
  - status in [200, 201]
  - duration < 1000
  - header Content-Type == application/json
  - jsonpath $.data[0].breed matches ^Abyss.+
captures: # captures can be used in later requests using the ${variable} syntax
  firstBreed: jsonpath $.data[0].breed
  secondBreed: jsonpath $.data[1].breed
```

To run this request, we'll go to the command-line:

```bash
$ nap run ./request-1.yml
Run finished in 297ms. 1/1 succeeded.
```

{: .highlight }
For the full request reference, containing all the options, see [File Types -> Request](/reference/file-types/requests)

## Routines

What makes Nap _FAST_? Routines make Nap _FAST_. Here's what a routine looks like in Nap:

```yml
# routine-1.yml
kind: routine
name: main routine
steps:
  - run: request-1.yml
```

We can run a routine the same way we would run a request:

```bash
$ nap run ./routine-1.yml
Run finished in 312ms. 1/1 succeeded.
```

Oops. That was actually slower. Let's change this up a bit. Let's start by adding a second request:

```yml
# request-2.yml
kind: request
name: Cat Facts - Succeeds
path: https://catfact.ninja/facts
asserts:
  - status == 200
  - duration < 1000
```

Simple enough. Now let's add a couple more routines. These will be redundant, but we're just illustrating a point here:

```yml
# subroutine-1.yml
kind: routine
name: example subroutine 1
steps:
  - run: request-1.yml
  - run: request-2.yml
```

```yml
# subroutine-2.yml
kind: routine
name: example subroutine 2
steps:
  - run: request-2.yml
  - run: request-1.yml
```

Finally, we'll update `routine-1.yml` to the following:

```yml
# routine-1.yml
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
Be sure not to reference a routine from itself, even indirectly. This will create an infinite loop, and your workload will stall!

You can view the full code for example on Github [here](https://github.com/davesheldon/nap/tree/main/examples/routines/basic).