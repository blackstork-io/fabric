---
title: blackstork/sqlite
weight: 20
type: docs
---

# `blackstork/sqlite` plugin

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/sqlite" = "=> v0.0.0-dev"
  }
}
```

## Data sources

The plugin has the following data sources available:

### `sqlite`

#### Configuration

The data source supports the following configuration parameters:

```hcl
config data sqlite {
    database_uri = <string>  # required
}
```

#### Usage

The data source supports the following parameters in the data blocks:

```hcl
data sqlite {
    sql_args = <list of dynamic>  # optional
    sql_query = <string>  # required
}
```