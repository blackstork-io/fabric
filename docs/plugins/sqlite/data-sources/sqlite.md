---
title: sqlite
plugin:
  name: blackstork/sqlite
  description: ""
  tags: []
  version: "v0.0.0-dev"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/sqlite/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/sqlite" "sqlite" "v0.0.0-dev" "sqlite" "data source" >}}

## Installation

To use `sqlite` data source, you must install the plugin `blackstork/sqlite`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/sqlite" = ">= v0.0.0-dev"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data sqlite {
    database_uri = <string>  # required
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data sqlite {
    sql_args = <list of dynamic>  # optional
    sql_query = <string>  # required
}
```