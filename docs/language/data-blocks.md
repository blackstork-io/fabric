---
title: Data Blocks
type: docs
weight: 50
---

# Data blocks

`data` blocks define data requirements for the template. In its signature, the block specifies the data source Fabric will fetch the data from.

```hcl
# Root-level definition of a data block
data <data-source-name> "<block-name>" {
  ...
}

document "foobar" {

  # In-document definition of a data block
  data <data-source-name> "<block-name>" {
    ...
  }

}
```

Both a data source name and a block name are required. The pair makes an unique identifier for the block.

The data blocks must be placed either on a root-level of the file or on a root-level of the document.

When Fabric starts rendering the template, the data sources for the data blocks are executed and the results are placed under the block names in the context (see [Context]({{< ref "content-blocks.md#context" >}}) for more details), available for queries.

## Supported arguments

The arguments provided in the block are either generic arguments or plugin-specific input parameters.

### Generic arguments

- `config`: (optional) a reference to a named configuration block for the data source or a content provider. If provided, it takes precedence over the default configuration. See [Data source configuration]({{< ref "configs.md#data-source-configuration" >}}) for the details.

### Plugin-specific Arguments

Plugin-specific arguments are defined by a plugin specification. See [Plugins]({{< ref "plugins.md" >}}) for the details on the supported arguments per plugin.

## Supported Nested Blocks

- `meta`: (optional) a block containing metadata for the block.
- `config`: (optional) an inline configuration for the block. If provided, it takes precedence over the `config` argument and default configuration for the plugin.

## References

See [References]({{< ref references.md >}}) for the details about referencing data blocks.

## Example

```hcl
config data csv {
  delimiter = ";"
}

data csv "events_a" {
  path = "/tmp/events-a.csv"
}

document "test-document" {

   data ref {
     base = data.csv.events_a
   }

   data csv "events_b" {
     config {
       delimiter = ",";
     }

     path = "/tmp/events-b.csv"
   }
}
```
