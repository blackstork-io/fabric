---
title: Fabric CLI
description: Fabric CLI is your gateway to Fabric's powerful features. Use `fabric` binary with commands like `install`, `data`, and `render`. Dive deeper with `fabric --help` to explore additional options and commands. Get started with Fabric CLI and unlock seamless document generation.
images:
  - 'images/diagram.png'
type: docs
weight: 20
---

# Fabric CLI

The command line interface to Fabric is `fabric` CLI tool.

The core Fabric commands are:

- `install` — installs all required plugins, listed in the [global configuration]({{< ref "language/configs.md#global-configuration" >}}). See [plugin installation docs]({{< ref "install.md#installing-plugins" >}}) for more details.
- `data` — executes the data block and prints out prettified JSON to standard output.
- `render` — renders the specified target (a document template) and prints out the result to standard output or to a file.

To get more details, run `fabric --help`:

```text
$ ./fabric --help
Usage:
  fabric [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  data        Execute a single data block
  help        Help about any command
  install     Install plugins
  render      Render the document

Flags:
      --color               enables colorizing the logs and diagnostics (if supported by the terminal and log format) (default true)
  -h, --help                help for fabric
      --log-format string   format of the logs (plain or json) (default "plain")
      --log-level string    logging level ('debug', 'info', 'warn', 'error') (default "info")
      --source-dir string   a path to a directory with *.fabric files (default ".")
  -v, --verbose             a shortcut to --log-level debug
      --version             version for fabric

Use "fabric [command] --help" for more information about a command.
```

## Source directory

Fabric loads `*.fabric` files from a source directory. By default, a source directory is the current directory (`.`) but it's possible to set `--source-dir` argument when running `fabric` to load files from a different location.

## Next step

Take a look at [the tutorial]({{< ref "tutorial.md" >}}) to see the commands in action.
