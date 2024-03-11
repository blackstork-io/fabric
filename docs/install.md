---
title: Install
type: docs
weight: 10
code_blocks_no_wrap: true
---

# Install

## Installing Fabric

### GitHub releases

Fabric binaries for Windows, macOS, and Linux are available at ["Releases"](https://github.com/blackstork-io/fabric/releases) page in the project's GitHub.

To install Fabric:

- **download a release archive**: choose and download Fabric release archive appropriate for your operating system (Windows, macOS/Darwin, or Linux) and architecture from ["Releases"](https://github.com/blackstork-io/fabric/releases) page;
- **unpack**: extract the contents of the downloaded release archive into a preferred directory;

For example, the steps for macOS (arm64) are:

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

Fabric relies on [the plugins]({{< ref "plugins.md" >}}) for implementing the integrations with various data sources, platforms, and services. Before the plugins can be used during the template rendering, they must be installed. Fabric's sub-command `install` can install the plugins automatically from the registry (`https://registry.blackstork.io`).

To install the plugins:

- **add all necessary plugins into the global configuration**: [the global configuration]({{< ref "language/configs.md#global-configuration" >}}) has a list of plugins dependencies in `plugin_versions` map. Add the plugins you would like to install in the map with a preferred version constraint.

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

The plugins are downloaded and installed in `./fabric` folder or in the location specified in `cache_dir` in [the global configuration]({{< ref "language/configs.md#global-configuration" >}}).

{{< hint note >}}

It's not necessary to install any plugins if you are only using built-in [data sources]({{< ref "plugins/builtin/_index.md#data-sources" >}}) and [content providers]({{< ref "plugins/builtin/_index.md#content-providers" >}}) in your templates

{{</ hint >}}

## Next step

That's it! You're now ready to use Fabric. Take a look at [the tutorial]({{< ref "tutorial.md" >}}) to see how to create and render the templates.
