---
title: blackstork/virustotal
weight: 20
type: docs
---

# `blackstork/virustotal` plugin

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/virustotal" = "=> v0.0.0-dev"
  }
}
```

## Data sources

The plugin has the following data sources available:

### `virustotal_api_usage`

#### Configuration

The data source supports the following configuration parameters:

```hcl
config data virustotal_api_usage {
    api_key = <string>  # required
}
```

#### Usage

The data source supports the following parameters in the data blocks:

```hcl
data virustotal_api_usage {
    end_date = <string>  # optional
    group_id = <string>  # optional
    start_date = <string>  # optional
    user_id = <string>  # optional
}
```