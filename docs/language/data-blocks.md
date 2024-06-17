---
title: Data Blocks
description: Learn how to use Fabric data blocks for defining data requirements for the templates.
type: docs
weight: 50
---
# Data blocks

`data` blocks define data requirements for the template. The data block represents a call to an
integration with provided parameters.

The block signature includes the name of the data source that will execute the data block.

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

Both a data source name and a block name are required, making an unique identifier for the block.

The data blocks must be placed either on a root-level of the file or on a root-level of the
document.

When Fabric renders the template, the data blocks are executed and the results are stored in the
context (see [Context]({{< ref context.md >}}) for more details), available for other blocks to use.

## Supported arguments

The arguments provided in the block are either generic arguments or data-source-specific arguments.

### Generic arguments

- `config`: (optional) a reference to a named configuration block for the data source. If provided,
  it takes precedence over the default configuration. See data source [configuration details]({{<
  ref "configs.md#block-configuration" >}}) for more information.

### Data source arguments

Data source arguments differ per data source. See the documentation for a specific data source (find
it in [supported data sources]({{< ref "data-sources.md" >}})) for the details on the supported
arguments.

## Supported nested blocks

- `meta`: (optional) a block containing metadata for the block.
- `config`: (optional) an inline configuration for the block. If provided, it takes precedence over
  the `config` argument and default configuration for the data source.
- `vars`: (optional) a block with variable definitions. See [Variables]({{< ref
  "context.md#variables" >}}) for the details.

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
       delimiter = ","
     }

     path = "/tmp/events-b.csv"
   }
}
```

## Next steps

See [Content Blocks]({{< ref "content-blocks.md" >}}) documentation to learn how to define content, like text paragraphs, tables, graphs and images, in the template.
