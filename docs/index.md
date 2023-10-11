---
layout: default
title: Home
nav_order: 1
---
{: .fs-10 .fw-300 }
# Lightning-Fast, Accurate, Open Source API Test Automation

{: .fs-6 .fw-300 }
Nap is a command-line tool that utilizes YAML files to test APIs.

![mpb](https://github.com/davesheldon/nap/assets/7782805/e7d70e5e-6ea3-4bb8-b909-e770ec9298f7)

{: .fs-6 .fw-300 }
Run large-scale workloads in a fraction of the time.

Nap super-charges your test automation workflow with parallel execution, breaking up test cases into groups that can be run in parallel.

{: .fs-6 .fw-300 }
Quickly and collaboratively test your HTTP APIs.

Nap's syntax is simple enough to quickly jot down during the design process. Powerful asserts make writing test cases a breeze. Nap's declarative YAML format means you can check your tests into source control to share them with your team.

{: .fs-6 .fw-300 }
Run locally or integrate with your favorite CI/CD tool. 

Nap compiles cross-platform to a single executable file. Run in Windows, Linux or Mac OS with ease. 

{: .fs-6 }
[Get Started](#getting-started){: .btn .btn-primary } [View on Github](https://github.com/davesheldon/nap){: .btn }

---

Nap is a command-line interface (CLI) for running HTTP requests using YAML files with a clear, concise syntax. You can write routines, requests and even limited scripts using your favorite text editor. You can then check those files into source control and intregate Nap with your CI/CD pipeline.

Browse the docs to learn more about how Nap can save you and your customers time and frustration with better, faster test automation.

## Why Nap?

Here are just a few of the reasons our users enjoy Nap:

- **_REMARKABLE SPEED._** Nap is able to break apart your test suite and run different parts of it at the same time. Even large-scale workloads only take about as long as your slowest scenario.
- **_SIMPLIFIED SYNTAX._** With Nap, there's no need to memorize a bunch of cURL flags or open a big fancy editor to design your tests. Each request is a single, compact file.
- **_VERSATILE EXPRESSIONS._** Write your tests in expressions that make sense at first glance. Asserts, variable captures and even explicit javascript are all at your disposal.

## Getting Started

There are two ways to get started using Nap.

### Download

You can download directly from the dist folder on Github, or use one of the quick download buttons:

[Download for Windows (x64) - v0.4.5](https://github.com/davesheldon/nap/releases/download/v0.4.5/nap.exe){: .btn .btn-primary} [Other Platforms (View All)](https://github.com/davesheldon/nap/tree/main/dist/){: .btn}

Once downloaded, add nap's location into your `$PATH` for convenience.

### Using Go

If you already have the Go language installed, simply run the command to install nap:

```bash
$ go install github.com/davesheldon/nap@latest
```

Once you have Nap installed, go read about [The Basics](/the-basics) to start writing your first test.
