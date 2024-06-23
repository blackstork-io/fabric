---
title: Install
description: Learn how to install Fabric and its plugins to streamline your templating workflow. Fabric binaries for Windows, macOS, and Linux are available at the project's GitHub releases page. Simply download the appropriate release archive for your operating system, unpack it, and you're ready to go.
type: docs
weight: 10
code_blocks_no_wrap: true
---

# Install

## Installing Fabric

### Homebrew

To install Fabric on macOS with [Homebrew](https://brew.sh/), run these commands:

```bash
# Install Fabric from the tap
brew install blackstork-io/tools/fabric

# Verify the version installed
fabric --version
```

It's recommended to use `blackstork-io/tools` tap when installing Fabric with Homebrew. The tap is
updated automatically with every release.

### Docker

The Docker images for Fabric are hosted in [Docker Hub](https://hub.docker.com/r/blackstorkio/fabric/tags).

To run Fabric as a Docker image use a full name `blackstorkio/fabric`:

```bash
docker run blackstorkio/fabric
```

### GitHub releases

Fabric binaries for Windows, macOS, and Linux are available at ["Releases"](https://github.com/blackstork-io/fabric/releases) page in the project's GitHub.

To install Fabric:

- **download a release archive**: choose and download Fabric release archive appropriate for your operating system (Windows, macOS/Darwin, or Linux) and architecture from ["Releases"](https://github.com/blackstork-io/fabric/releases) page;
- **unpack**: extract the contents of the downloaded release archive into a preferred directory;

For example, the steps to install the latest Fabric release from GitHub on macOS (arm64) are:

```bash
# Create a folder
mkdir fabric-bin

# Download the latest release of Fabric
wget https://github.com/blackstork-io/fabric/releases/latest/download/fabric_darwin_arm64.tar.gz -O ./fabric_darwin_arm64.tar.gz

# Unpack Fabric release archive into `fabric-bin` folder
tar -xvzf ./fabric_darwin_arm64.tar.gz -C ./fabric-bin

# Verify that `fabric` runs
./fabric-bin/fabric --help
```

## Installing plugins

Fabric uses [plugins]({{< ref "plugins.md" >}}) for integrations with various data sources, platforms, and services.

Before the plugins can be used for template rendering, they must be installed. Fabric has a sub-command `install` that automatically installs all required plugins from the registry (`https://registry.blackstork.io` is the root endpoint for the registry).

To install the plugins:

- **add necessary plugins to the list in global configuration**: [the global configuration]({{< ref "language/configs.md#global-configuration" >}}) has a list of plugins dependencies in `plugin_versions` map. Add the plugins to install in the map with a preferred version constraint.

  ```hcl
  fabric {
    plugin_versions = {
      "blackstork/openai" = ">= 0.0.1",
      "blackstork/elastic" = ">= 0.0.1",
    }
  }
  ```

- **install the plugins**: run `install` sub-command to install the plugins. For example:

  ```text
  $ ./fabric install
  Mar 11 19:20:09.085 INF Searching plugin name=blackstork/elastic constraints=">=v0.0.1"
  Mar 11 19:20:09.522 INF Installing plugin name=blackstork/elastic version=0.4.0
  Mar 11 19:20:10.769 INF Searching plugin name=blackstork/openai constraints=">=v0.0.1"
  Mar 11 19:20:10.787 INF Installing plugin name=blackstork/openai version=0.4.0
  $
  ```

Fabric downloads and installs plugins in `./.fabric` folder, or in the location specified in `cache_dir` in [the global configuration]({{< ref "language/configs.md#global-configuration" >}}).

{{< hint note >}}
There is no need to install plugins if you are only using resources from a [built-in plugin]({{< ref "plugins/builtin/_index.md" >}}) in the templates.
{{</ hint >}}

## Next step

That's it! You're now ready to use Fabric. Take a look at [Tutorial]({{< ref "tutorial.md" >}}) to see how to create and render the templates.
