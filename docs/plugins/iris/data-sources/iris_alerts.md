---
title: "`iris_alerts` data source"
plugin:
  name: blackstork/iris
  description: "Retrieve alerts from Iris API"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/iris/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/iris" "iris" "v0.4.2" "iris_alerts" "data source" >}}

## Description
Retrieve alerts from Iris API

## Installation

To use `iris_alerts` data source, you must install the plugin `blackstork/iris`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/iris" = ">= v0.4.2"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration arguments:

```hcl
config data iris_alerts {
  # Iris API url
  #
  # Required string.
  # Must be non-empty
  #
  # For example:
  api_url = "some string"

  # Iris API Key
  #
  # Required string.
  # Must be non-empty
  #
  # For example:
  api_key = "some string"

  # Enable/disable insecure TLS
  #
  # Optional bool.
  # Default value:
  insecure = false
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data iris_alerts {
  # List of Alert IDs
  #
  # Optional list of number.
  # Default value:
  alert_ids = null

  # Alert Source
  #
  # Optional string.
  # Default value:
  alert_source = null

  # List of tags
  #
  # Optional list of string.
  # Default value:
  tags = null

  # Case ID
  #
  # Optional number.
  # Default value:
  case_id = null

  # Alert Customer ID
  #
  # Optional number.
  # Default value:
  customer_id = null

  # Alert Owner ID
  #
  # Optional number.
  # Default value:
  owner_id = null

  # Alert Severity ID
  #
  # Optional number.
  # Default value:
  severity_id = null

  # Alert Classification ID
  #
  # Optional number.
  # Default value:
  classification_id = null

  # Alert State ID
  #
  # Optional number.
  # Default value:
  status_id = null

  # Alert Date - lower boundary
  #
  # Optional string.
  # Default value:
  alert_start_date = null

  # Alert Date - higher boundary
  #
  # Optional string.
  # Default value:
  alert_end_date = null

  # Sort order
  #
  # Optional string.
  # Must be one of: "desc", "asc"
  # Default value:
  sort = "desc"

  # Size limit to retrieve
  #
  # Optional number.
  # Must be >= 0
  # Default value:
  size = 0
}
```