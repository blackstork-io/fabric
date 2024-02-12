<div align="center">

<img src=".assets/fabric.svg" alt="fabric-logo" width="250px"/>
<br/>
<br/>

Codifying and automating mission-critical communications with standardized and reusable templates.

#### [Releases](https://github.com/blackstork-io/fabric/releases) | [Docs](https://blackstork.io/fabric/docs/) | [Slack](https://fabric-community.slack.com/)

![GitHub Repo stars](https://img.shields.io/github/stars/blackstork-io/fabric?style=social)
![GitHub Release](https://img.shields.io/github/v/release/blackstork-io/fabric)
[![Join Slack](https://img.shields.io/badge/slack-join-8F87F7)](https://fabric-community.slack.com/)

</div>

> [!NOTE]  
> Fabric is currently in the early stages of development, and there may be some issues. We welcome your feedback, so if you have any suggestions, ideas, or encounter bugs, please share them with us in our [Fabric Community slack](https://fabric-community.slack.com/).

Fabric is an open-source configuration language and a CLI tool that enables the codification and automation of the content generation process.

<div align="center">
    <img src=".assets/diagram.svg" alt="fabric-diagram" width="600px"/>
</div>

Fabric produces Markdown documents from the templates that declaratively define data requirements and content structure. The templates are written in Fabric Configuration Language and consist of reusable blocks, powered by plugins.
Data blocks fetch data from various external sources -- data stores, security solutions, and platforms. The content blocks render the template into Markdown document.

See [Documentation](https://blackstork.io/fabric/docs/) for more details on the Fabric language and Fabric CLI.


# Installation

To get started with Fabric, follow these simple steps for installation across various operating systems:

- **download release archives**: choose and download the appropriate release for your operating system (Windows, macOS, or Linux) and architecture (32-bit or 64-bit) in ["Releases" section](https://github.com/blackstork-io/fabric/releases);
- **unpack the archives**: extract the contents of the downloaded archive to a preferred directory;
- **run the binary**: run `fabric` binary from the command line to launch Fabric.

That's it! You're now ready to use Fabric. For more details on usage and configuration options, refer to the "Usage" paragraph below or [Fabric CLI](https://blackstork.io/fabric/docs/cli) documentation.


# Usage

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
      --color                enables colorizing the logs and diagnostics (if supported by the terminal and log format) (default true)
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

# Documentation

Visit [https://blackstork.io/fabric/docs/](https://blackstork.io/fabric/docs/) for full documentation.

# Security

Please report any suspected security vulnerabilities through GitHub's [security advisory reporting](https://github.com/blackstork-io/fabric/security/advisories/new). We take all legitimate reports seriously and will thoroughly investigate.

We kindly request that you talk to us before making any public disclosures. This ensures that no excessive information is revealed before a patch is ready and users have sufficient time to upgrade.

# License

Fabric is licensed under Apache-2.0 license. See the [LICENSE](LICENSE) file for the details.
