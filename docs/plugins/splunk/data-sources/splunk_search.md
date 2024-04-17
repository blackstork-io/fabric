---
title: splunk_search
plugin:
  name: blackstork/splunk
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/splunk/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/splunk" "splunk" "v0.4.1" "splunk_search" "data source" >}}

## Installation

To use `splunk_search` data source, you must install the plugin `blackstork/splunk`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/splunk" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data splunk_search {
  # Required string. For example:
  auth_token = "some string"

  # Optional string. Default value:
  host = null

  # Optional string. Default value:
  deployment_name = null
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data splunk_search {
  # Required string. For example:
  search_query = "some string"

  # Optional number. Default value:
  max_count = null

  # Optional number. Default value:
  status_buckets = null

  # Optional list of string. Default value:
  rf = null

  # Optional string. Default value:
  earliest_time = null

  # Optional string. Default value:
  latest_time = null
}
```