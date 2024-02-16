---
title: blackstork/postgresql
weight: 20
type: docs
---

# `blackstork/postgresql` plugin

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/postgresql" = "=> v0.0.0-dev"
  }
}
```

## Data sources

The plugin has the following data sources available:

### `postgresql`

#### Configuration

The data source supports the following configuration parameters:

```hcl
config data postgresql {
    database_url = <string>  # required
}
```

#### Usage

The data source supports the following parameters in the data blocks:

```hcl
data postgresql {
    sql_args = <list of dynamic>  # optional
    sql_query = <string>  # required
}
```