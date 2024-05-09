---
title: elastic_security_cases
plugin:
  name: blackstork/elastic
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/elastic/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/elastic" "elastic" "v0.4.1" "elastic_security_cases" "data source" >}}

## Installation

To use `elastic_security_cases` data source, you must install the plugin `blackstork/elastic`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/elastic" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data elastic_security_cases {
  # Required string.
  # For example:
  kibana_endpoint_url = "some string"

  # Optional string.
  # Default value:
  api_key_str = null

  # Optional list of string.
  # Default value:
  api_key = null
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data elastic_security_cases {
  # Optional string.
  # Default value:
  space_id = null

  # Optional list of string.
  # Default value:
  assignees = null

  # Optional string.
  # Default value:
  default_search_operator = null

  # Optional string.
  # Default value:
  from = null

  # Optional list of string.
  # Default value:
  owner = null

  # Optional list of string.
  # Default value:
  reporters = null

  # Optional string.
  # Default value:
  search = null

  # Optional list of string.
  # Default value:
  search_fields = null

  # Optional string.
  # Default value:
  severity = null

  # Optional string.
  # Default value:
  sort_field = null

  # Optional string.
  # Default value:
  sort_order = null

  # Optional string.
  # Default value:
  status = null

  # Optional list of string.
  # Default value:
  tags = null

  # Optional string.
  # Default value:
  to = null

  # Optional number.
  # Default value:
  size = null
}
```