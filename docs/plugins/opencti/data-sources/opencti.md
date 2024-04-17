---
title: opencti
plugin:
  name: blackstork/opencti
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/opencti/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/opencti" "opencti" "v0.4.1" "opencti" "data source" >}}

## Installation

To use `opencti` data source, you must install the plugin `blackstork/opencti`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/opencti" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data opencti {
  # Required string. For example:
  graphql_url = "some string"

  # Optional string. Default value:
  auth_token = null
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data opencti {
  # Required string. For example:
  graphql_query = "some string"
}
```