---
title: microsoft_sentinel_incidents
plugin:
  name: blackstork/microsoft_sentinel
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/sentinel/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/microsoft_sentinel" "microsoft_sentinel" "v0.4.1" "microsoft_sentinel_incidents" "data source" >}}

## Installation

To use `microsoft_sentinel_incidents` data source, you must install the plugin `blackstork/microsoft_sentinel`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/microsoft_sentinel" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data microsoft_sentinel_incidents {
    resource_group_name = <string>  # required
    subscription_id = <string>  # required
    workspace_name = <string>  # required
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data microsoft_sentinel_incidents {
    filter = <string>  # optional
    limit = <number>  # optional
    order_by = <string>  # optional
}
```