---
title: "`elasticsearch` data source"
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

{{< plugin-resource-header "blackstork/elastic" "elastic" "v0.4.1" "elasticsearch" "data source" >}}

## Installation

To use `elasticsearch` data source, you must install the plugin `blackstork/elastic`.

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

The data source supports the following configuration arguments:

```hcl
config data elasticsearch {
  # Optional string.
  # Default value:
  base_url = null

  # Optional string.
  # Default value:
  cloud_id = null

  # Optional string.
  # Default value:
  api_key_str = null

  # Optional list of string.
  # Default value:
  api_key = null

  # Optional string.
  # Default value:
  basic_auth_username = null

  # Optional string.
  # Default value:
  basic_auth_password = null

  # Optional string.
  # Default value:
  bearer_auth = null

  # Optional string.
  # Default value:
  ca_certs = null
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data elasticsearch {
  # Required string.
  # For example:
  index = "some string"

  # Optional string.
  # Default value:
  id = null

  # Optional string.
  # Default value:
  query_string = null

  # Optional map of any single type.
  # Default value:
  query = null

  # Optional map of any single type.
  # Default value:
  aggs = null

  # Optional bool.
  # Default value:
  only_hits = null

  # Optional list of string.
  # Default value:
  fields = null

  # Optional number.
  # Default value:
  size = null
}
```