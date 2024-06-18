---
title: "`virustotal_api_usage` data source"
plugin:
  name: blackstork/virustotal
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/virustotal/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/virustotal" "virustotal" "v0.4.1" "virustotal_api_usage" "data source" >}}

## Installation

To use `virustotal_api_usage` data source, you must install the plugin `blackstork/virustotal`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/virustotal" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration arguments:

```hcl
config data virustotal_api_usage {
  # Required string.
  # For example:
  api_key = "some string"
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data virustotal_api_usage {
  # Optional string.
  # Default value:
  user_id = null

  # Optional string.
  # Default value:
  group_id = null

  # Optional string.
  # Default value:
  start_date = null

  # Optional string.
  # Default value:
  end_date = null
}
```