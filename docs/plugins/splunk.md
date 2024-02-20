---
title: blackstork/splunk
weight: 20
type: docs
---

# `blackstork/splunk` plugin

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/splunk" = "=> v0.0.0-dev"
  }
}
```

## Data sources

The plugin has the following data sources available:

### `splunk_search`

#### Configuration

The data source supports the following configuration parameters:

```hcl
config data splunk_search {
    auth_token = <string>  # required
    deployment_name = <string>  # optional
    host = <string>  # optional
}
```

#### Usage

The data source supports the following parameters in the data blocks:

```hcl
data splunk_search {
    earliest_time = <string>  # optional
    latest_time = <string>  # optional
    max_count = <number>  # optional
    rf = <list of string>  # optional
    search_query = <string>  # required
    status_buckets = <number>  # optional
}
```