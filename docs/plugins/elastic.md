---
title: blackstork/elastic
weight: 20
type: docs
---

# `blackstork/elastic` plugin

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/elastic" = "=> v0.0.0-dev"
  }
}
```

## Data sources

The plugin has the following data sources available:

### `elasticsearch`

#### Configuration

The data source supports the following configuration parameters:

```hcl
config data elasticsearch {
    api_key = <list of string>  # optional
    api_key_str = <string>  # optional
    base_url = <string>  # optional
    basic_auth_password = <string>  # optional
    basic_auth_username = <string>  # optional
    bearer_auth = <string>  # optional
    ca_certs = <string>  # optional
    cloud_id = <string>  # optional
}
```

#### Usage

The data source supports the following parameters in the data blocks:

```hcl
data elasticsearch {
    aggs = <map of dynamic>  # optional
    fields = <list of string>  # optional
    id = <string>  # optional
    index = <string>  # required
    only_hits = <bool>  # optional
    query = <map of dynamic>  # optional
    query_string = <string>  # optional
    size = <number>  # optional
}
```