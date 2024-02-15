---
title: blackstork/graphql
weight: 20
type: docs
---

# `blackstork/graphql` plugin

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/graphql" = "=> v0.0.0-dev"
  }
}
```

## Data sources

The plugin has the following data sources available:

### `graphql`

#### Configuration

The data source supports the following configuration parameters:

```hcl
config data graphql {
    auth_token = <string>  # optional
    url = <string>  # required
}
```

#### Usage

The data source supports the following parameters in the data blocks:

```hcl
data graphql {
    query = <string>  # required
}
```