---
title: "`graphql` data source"
plugin:
  name: blackstork/graphql
  description: ""
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/graphql/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/graphql" "graphql" "v0.4.2" "graphql" "data source" >}}

## Installation

To use `graphql` data source, you must install the plugin `blackstork/graphql`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/graphql" = ">= v0.4.2"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration arguments:

```hcl
config data graphql {
  # Required string.
  #
  # For example:
  url = "some string"

  # Optional string.
  # Default value:
  auth_token = null
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data graphql {
  # Required string.
  #
  # For example:
  query = "some string"
}
```