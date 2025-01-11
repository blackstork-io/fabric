---
title: "`iris_cases` data source"
plugin:
  name: blackstork/iris
  description: "Retrieve cases from Iris API"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/iris/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/iris" "iris" "v0.4.2" "iris_cases" "data source" >}}

## Description
Retrieve cases from Iris API

## Installation

To use `iris_cases` data source, you must install the plugin `blackstork/iris`.

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
config data iris_cases {
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
data iris_cases {
  # List of Case IDs
  #
  # Optional list of number.
  # Default value:
  case_ids = null

  # Case Customer ID
  #
  # Optional number.
  # Default value:
  customer_id = null

  # Case Owner ID
  #
  # Optional number.
  # Default value:
  owner_id = null

  # Case Severity ID
  #
  # Optional number.
  # Default value:
  severity_id = null

  # Case State ID
  #
  # Optional number.
  # Default value:
  state_id = null

  # Case SOC ID
  #
  # Optional string.
  # Default value:
  soc_id = null

  # Case opening date - lower boundary
  #
  # Optional string.
  # Default value:
  start_open_date = null

  # Case opening date - higher boundary
  #
  # Optional string.
  # Default value:
  end_open_date = null

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