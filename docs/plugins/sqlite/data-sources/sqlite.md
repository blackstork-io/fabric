---
title: "`sqlite` data source"
plugin:
  name: blackstork/sqlite
  description: ""
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/sqlite/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/sqlite" "sqlite" "v0.4.2" "sqlite" "data source" >}}

## Installation

To use `sqlite` data source, you must install the plugin `blackstork/sqlite`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/sqlite" = ">= v0.4.2"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration arguments:

```hcl
config data sqlite {
  # Required string.
  # For example:
  database_uri = "some string"
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data sqlite {
  # Required string.
  # For example:
  sql_query = "some string"

  # Optional list of any single type.
  # Default value:
  sql_args = null
}
```