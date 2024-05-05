---
title: graphql
plugin:
  name: blackstork/graphql
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/graphql/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/graphql" "graphql" "v0.4.1" "graphql" "data source" >}}

## Installation

To use `graphql` data source, you must install the plugin `blackstork/graphql`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/graphql" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data graphql {
  # API endpoint to perform GraphQL queries against
  #
  # Required string.
  # For example:
  url = "https://example.com/graphql"

  # Token to be sent to the server as "Authorization: Bearer" header.
  # Empty or null tokens are not sent.
  #
  # Optional string.
  # For example:
  # auth_token = "<token>"
  # 
  # Default value:
  auth_token = null
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data graphql {
  # GraphQL query
  #
  # Required string.
  # For example:
  query = "query{user{id, name}}"
}
```