---
title: blackstork/elasticsearch
weight: 20
type: docs
---

# `blackstork/elasticsearch` plugin

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/elasticsearch" = "=> v0.0.0-dev"
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
    base_url = <string>  # required
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
    fields = <list of string>  # optional
    id = <string>  # optional
    index = <string>  # required
    query = <map of dynamic>  # optional
    query_string = <string>  # optional
}
```