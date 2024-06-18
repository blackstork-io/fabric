---
title: "`microsoft_sentinel_incidents` data source"
plugin:
  name: blackstork/microsoft
  description: "The `microsoft_sentinel_incidents` data source fetches incidents from Microsoft Sentinel"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/microsoft/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/microsoft" "microsoft" "v0.4.1" "microsoft_sentinel_incidents" "data source" >}}

## Description
The `microsoft_sentinel_incidents` data source fetches incidents from Microsoft Sentinel.

## Installation

To use `microsoft_sentinel_incidents` data source, you must install the plugin `blackstork/microsoft`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/microsoft" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration arguments:

```hcl
config data microsoft_sentinel_incidents {
  # The Azure client ID
  #
  # Required string.
  # For example:
  client_id = "some string"

  # The Azure client secret
  #
  # Required string.
  # For example:
  client_secret = "some string"

  # The Azure tenant ID
  #
  # Required string.
  # For example:
  tenant_id = "some string"

  # The Azure subscription ID
  #
  # Required string.
  # For example:
  subscription_id = "some string"

  # The Azure resource group name
  #
  # Required string.
  # For example:
  resource_group_name = "some string"

  # The Azure workspace name
  #
  # Required string.
  # For example:
  workspace_name = "some string"
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data microsoft_sentinel_incidents {
  # The filter expression
  #
  # Optional string.
  # Default value:
  filter = null

  # The maximum number of incidents to return
  #
  # Optional number.
  # Default value:
  limit = null

  # The order by expression
  #
  # Optional string.
  # Default value:
  order_by = null
}
```