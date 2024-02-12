---
title: Fabric CLI
type: docs
weight: 4
---

# Fabric CLI

The command line interface to Fabric is `fabric` CLI tool. It supports two subcommands:

- `fabric data` — executes the data block and prints out prettified JSON to stdout.
- `fabric render` — renders the specified target (a document template) into Markdown and outputs the result to stdout or to a file.

```bash
$ fabric --help

Usage:
  fabric [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  data        Execute a single data block
  help        Help about any command
  render      Render the document

Flags:
      --color                enables colorizing the logs and diagnostics (if supported by the terminal
                             and log format) (default true)
  -h, --help                 help for fabric
      --log-format string    format of the logs (plain or json) (default "plain")
      --log-level string     logging level ('debug', 'info', 'warn', 'error') (default "info")
      --plugins-dir string   override for plugins dir from fabric configuration (required)
      --source-dir string    a path to a directory with *.fabric files (default ".")
  -v, --verbose              a shortcut to --log-level debug
      --version              version for fabric

Use "fabric [command] --help" for more information about a command.
```

## Source Directory

Fabric loads `*.fabric` files from a source directory. By default, a source directory is the current directory  (`.`), where `fabric` is executed. To provide a different location, use `--source-dir` argument when running `fabric`.

## Plugins Directory

For now, `fabric` expects the plugins directory path to be provided with `--plugins-dir` argument. This will change in the future, when `fabric` will perform plugin discovery automatically.

