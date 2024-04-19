---
title: sqlite
plugin:
  name: blackstork/sqlite
  description: "Produces query results from Sqlite"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/sqlite/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/sqlite" "sqlite" "v0.4.1" "sqlite" "data source" >}}

## Description
Produces query results from Sqlite

## Installation

To use `sqlite` data source, you must install the plugin `blackstork/sqlite`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/sqlite" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data sqlite {
  # Required string. For example:
  database_uri = "file:test.db"
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data sqlite {
  # Required string. For example:
  sql_query = "SELECT * FROM example WHERE id=$1 OR age=$2"

  # Values for the prepared statement
  #
  # For example:
  # sql_args = [42, 24]
  #
  # Optional list of any single type. Default value:
  sql_args = null
}
```