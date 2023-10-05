---
layout: default
title: Home
nav_order: 1
---
{: .fs-10 .fw-300 }
# Achieve _FAST_ HTTP test automation without all the fuss.

{: .fs-6 .fw-300 }
Nap super-charges your test automation workflow with parallel execution, powerful asserts and the flexibility of file-based storage.

![mpb](https://github.com/davesheldon/nap/assets/7782805/e7d70e5e-6ea3-4bb8-b909-e770ec9298f7)


{: .fs-6 }
[Get Started](#getting-started){: .btn .btn-primary } [View on Github](https://github.com/davesheldon/nap){: .btn }

---

Nap is a command-line interface (CLI) for running HTTP requests using YAML files with a clear, concise syntax. You can write routines, requests and even limited scripts using your favorite text editor. You can then check those files into source control and intregate Nap with your CI/CD pipeline.

Browse the docs to learn more about how Nap can save you and your customers time and frustration with better, faster test automation.

## Why Nap?

Here are just a few of the reasons our users enjoy Nap:

- **_LUDICROUS SPEED._** Nap is able to break apart your test suite and run different parts of it at the same time. So your entire run will only take about as long as your slowest scenario.
- **_EASY SYNTAX._** With Nap, there's no need to memorize a bunch of cURL flags or open a big fancy editor to design your tests. Each request is a single, compact file.
- **_POWERFUL EXPRESSIONS._** Write your tests in expressions that make sense at first glance. Asserts, variable captures and even explicit javascript are all at your disposal.

## Getting Started

There are two ways to get started using Nap.

### Download

You can download directly from the dist folder on Github, or use one of the quick download buttons:

[v0.4.0 - Windows (x64)](https://github.com/davesheldon/nap/releases/download/v0.4.0/nap.exe){: .btn}

Once downloaded, copy `nap.exe`'s location into your `$PATH` for convenience.

### Using Go

If you already have the Go language installed, simply run the command to install nap:

```bash
$ go install github.com/davesheldon/nap@latest
```

Once you have Nap installed, go read about [The Basics](/the-basics) to start writing your first test.
