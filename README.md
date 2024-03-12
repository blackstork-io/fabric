<div align="center">

<img src=".github/fabric.svg" alt="Fabric logo" width="250px"/>
<br/>
<br/>

Codifying and automating mission-critical communications with standardized and reusable templates.

[Releases](https://github.com/blackstork-io/fabric/releases) | [Docs](https://blackstork.io/fabric/docs/) | [Slack](https://fabric-community.slack.com/)

![GitHub Repository stars](https://img.shields.io/github/stars/blackstork-io/fabric?style=social)
![GitHub Release](https://img.shields.io/github/v/release/blackstork-io/fabric)
[![Join Slack](https://img.shields.io/badge/slack-join-8F87F7)](https://fabric-community.slack.com/)

</div>

**Fabric** is an open-source Command-Line Interface (CLI) tool with a configuration language designed to encode and automate content generation for cyber security and compliance.

<div align="center">
    <img src=".github/diagram.png" alt="The diagram illustrates a Fabric template and the corresponding rendered document" width="700px"/>
</div>

Fabric generates documents based on modular templates written in [Fabric Configuration Language](https://blackstork.io/fabric/docs/language/) (FCL). Template authors can delineate data requirements and content structures within the template, significantly reducing the manual effort associated with data consolidation and improving re-usability.

The plugin ecosystem expands Fabric's capabilities by offering integrations with various platforms, data stores, and security solutions, including [JSON/CSV files](https://blackstork.io/fabric/docs/plugins/builtin/), [Elastic](https://blackstork.io/fabric/docs/plugins/elastic/), [OpenCTI](https://blackstork.io/fabric/docs/plugins/opencti/), [Splunk Cloud](https://blackstork.io/fabric/docs/plugins/splunk/), [GitHub](https://blackstork.io/fabric/docs/plugins/github/), and more. A comprehensive list of supported plugins is available in [the documentation](https://blackstork.io/fabric/docs/plugins/).

To facilitate flexible content generation, Fabric content providers leverage powerful [Go templates](https://pkg.go.dev/text/template) and incorporate capabilities such as generative AI/LLMs (OpenAI API, Azure OpenAI, etc).

Refer to [the documentation](https://blackstork.io/fabric/docs/) for in-depth details about Fabric.

> [!NOTE]  
> Fabric is currently in the early stages of development, and there may be some issues. If you have any suggestions, ideas, or bug reports, please share them in [Fabric Community slack](https://fabric-community.slack.com/).

# Installation

## Installing Fabric

Fabric binaries for Windows, macOS, and Linux are available at ["Releases"](https://github.com/blackstork-io/fabric/releases) page in the project's GitHub.

To install Fabric:

- **download a release archive**: choose and download Fabric release archive appropriate for your operating system (Windows, macOS/Darwin, or Linux) and architecture from ["Releases"](https://github.com/blackstork-io/fabric/releases) page;
- **unpack**: extract the contents of the downloaded release archive into a preferred directory.

For example, the steps for macOS (arm64) are:

```bash
# Create a folder
mkdir fabric-bin

# Download the latest release of Fabric
wget https://github.com/blackstork-io/fabric/releases/latest/download/fabric_darwin_arm64.tar.gz \
    -O ./fabric_darwin_arm64.tar.gz

# Unpack Fabric release archive into `fabric-bin` folder
tar -xvzf ./fabric_darwin_arm64.tar.gz -C ./fabric-bin

# Verify that `fabric` runs
./fabric-bin/fabric --help
```

## Installing Fabric plugins

Fabric relies on [the plugins](https://blackstork.io/fabric/docs/plugins/) for implementing the integrations with various data sources, platforms, and services.

Before the plugins can be used during the template rendering, they must be installed. Fabric's sub-command `install` can install the plugins automatically from the registry.

To install the plugins:

- **add all necessary plugins into the global configuration**: [the global configuration](https://blackstork.io/fabric/docs/language/configs/#global-configuration) has a list of plugins dependencies in `plugin_versions` map. Add the plugins you want to install in the map with a preferred version constraint.

  ```hcl
  fabric {
    plugin_versions = {
      "blackstork/openai" = ">= 0.0.1",
      "blackstork/elastic" = ">= 0.0.1",
    }
  }
  ```

- **install the plugins**: run `install` sub-command to install the plugins. For example:

  ```bash
  $ ./fabric install
  Mar 11 19:20:09.085 INF Searching plugin name=blackstork/elastic constraints=">=v0.0.1"
  Mar 11 19:20:09.522 INF Installing plugin name=blackstork/elastic version=0.4.0
  Mar 11 19:20:10.769 INF Searching plugin name=blackstork/openai constraints=">=v0.0.1"
  Mar 11 19:20:10.787 INF Installing plugin name=blackstork/openai version=0.4.0
  $
  ```

The plugins are downloaded and installed in `./fabric` folder.

> [!NOTE]  
> It's not necessary to install any plugins if you are only using built-in [data sources](https://blackstork.io/fabric/docs/plugins/builtin/#data-sources) and [content providers](https://blackstork.io/fabric/docs/plugins/builtin/#content-providers) in your templates

# Usage

The command line interface to Fabric is `fabric` CLI tool.

The core sub-commands are:

- `install` — installs all required plugins.
- `data` — executes the data block and prints out prettified JSON to standard output.
- `render` — renders the specified target (a document template) and prints out the result to standard output or a file.

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

# Templates

You can find free Fabric templates in [Fabric Templates](https://github.com/blackstork-io/fabric-templates) repository.

# Documentation

Visit [https://blackstork.io/fabric/docs/](https://blackstork.io/fabric/docs/) for full documentation.

# Security

If you suspect any vulnerabilities within Fabric, please promptly report them using GitHub's [security advisory reporting](https://github.com/blackstork-io/fabric/security/advisories/new). We treat every report with the utmost seriousness and commit to conducting a thorough investigation.

We kindly request that you converse with us before making any public disclosures. This precautionary step ensures that no excessive information is divulged prematurely, allowing us to prepare a patch. Additionally, it gives users enough time to upgrade and enhance their system security.

# License

Fabric is licensed under Apache-2.0 license. See the [LICENSE](LICENSE) file for the details.
