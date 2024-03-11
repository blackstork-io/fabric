---
title: splunk_search
plugin:
  name: blackstork/splunk
  description: ""
  tags: []
  version: "v0.4.0"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/splunk/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/splunk" "splunk" "v0.4.0" "splunk_search" "data source" >}}

## Installation

To use `splunk_search` data source, you must install the plugin `blackstork/splunk`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/splunk" = ">= v0.4.0"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data splunk_search {
    auth_token = <string>  # required
    deployment_name = <string>  # optional
    host = <string>  # optional
}
```

## Usage

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