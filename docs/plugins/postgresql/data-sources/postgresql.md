---
title: postgresql
plugin:
  name: blackstork/postgresql
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/postgresql/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/postgresql" "postgresql" "v0.4.1" "postgresql" "data source" >}}

## Installation

To use `postgresql` data source, you must install the plugin `blackstork/postgresql`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/postgresql" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data postgresql {
    database_url = <string>  # required
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data postgresql {
    sql_args = <list of dynamic>  # optional
    sql_query = <string>  # required
}
```