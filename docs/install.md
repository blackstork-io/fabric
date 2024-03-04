---
title: Install
type: docs
weight: 10
---

# Install

## GitHub releases

The archives of compiled binaries for Fabric and Fabric plugins are available for Windows, macOS, and Linux at the ["Releases"](https://github.com/blackstork-io/fabric/releases) page.

To get started with Fabric:

- **download release archives**: choose and download the appropriate Fabric and Fabric plugin releases for your operating system (Windows, macOS/Darwin, or Linux) and architecture from ["Releases" section](https://github.com/blackstork-io/fabric/releases);
- **unpack the archives**: extract the contents of the downloaded archives to a preferred directory;

For example, the steps for macOS (arm64) are:

```bash
# Create a folder
mkdir fabric-bin

# Download the latest release of Fabric
wget https://github.com/blackstork-io/fabric/releases/latest/download/fabric_darwin_arm64.tar.gz -O ./fabric_darwin_arm64.tar.gz

# Download the latest release of Fabric plugins
wget https://github.com/blackstork-io/fabric/releases/latest/download/plugins_darwin_arm64.tar.gz -O ./plugins_darwin_arm64.tar.gz

# Unpack Fabric release archive into `fabric-bin` folder
tar -xvzf ./fabric_darwin_arm64.tar.gz -C ./fabric-bin

# Unpack Fabric plugins release archive into `fabric-bin` folder
tar -xvzf ./plugins_darwin_arm64.tar.gz -C ./fabric-bin

# Verify that `fabric` runs
./fabric-bin/fabric --help
```

That's it! You're now ready to use Fabric. For more details on usage and configuration options, refer [Fabric CLI]({{< ref "cli.md" >}}) documentation.
