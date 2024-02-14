---
title: Install
type: docs
weight: 2
---

# Install

## GitHub releases

The archives of compiled binaries for Fabric are available for Windows, macOS, and Linux at the project's ["Releases"](https://github.com/blackstork-io/fabric/releases) page.

To get started with Fabric:

- **download release archives**: choose and download the appropriate release for your operating system (Windows, macOS/Darwin, or Linux) and architecture in ["Releases" section](https://github.com/blackstork-io/fabric/releases);
- **unpack the archives**: extract the contents of the downloaded archive to a preferred directory;

For example, on macOS you can do:

```bash
wget https://github.com/blackstork-io/fabric/releases/download/v0.3.0/vale_0.3.0_Linux_64-bit.tar.gz
mkdir fabric && tar -xvzf vale_2.28.0_Linux_64-bit.tar.gz -C bin
export PATH=./bin:"$PATH"
```

That's it! You're now ready to use Fabric. For more details on usage and configuration options, refer [Fabric CLI]({{< ref "cli.md" >}}) documentation.
