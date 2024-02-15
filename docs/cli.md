---
title: Fabric CLI
type: docs
weight: 4
---

# Fabric CLI

The command line interface to Fabric is `fabric` CLI tool. It supports two sub-commands:

- `fabric data` — executes the data block and prints out prettified JSON to standard output
- `fabric render` — renders the specified target (a document template) into Markdown and outputs the result to standard output or to a file.

```text
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
      --plugins-dir string   override for plugins dir from fabric configuration
      --source-dir string    a path to a directory with *.fabric files (default ".")
  -v, --verbose              a shortcut to --log-level debug
      --version              version for fabric

Use "fabric [command] --help" for more information about a command.
```

## Source directory

Fabric loads `*.fabric` files from a source directory. By default, a source directory is the current directory  (`.`). To load Fabric files from a different location, set `--source-dir` argument when running `fabric`.
