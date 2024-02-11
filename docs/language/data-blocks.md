---
title: Data Blocks
type: docs
weight: 50
---

# Data Blocks

Blocks of type `data` define data requirements for the document. The data blocks are executed by data plugins.

```hcl
data <plugin-name> "<block-name>" {
  ...
}

document "foobar" {

  data <plugin-name> "<block-name>" {
    ...
  }

}
```

The plugin name and the block name are required, and, together with the block type `data` make a unique identifier for the block. The `data` blocks must be defined on the root-level of the file or of the document.

The data, represented by the block, is accessible under `data.<plugin-name>.<block-name>` path in the context (see [Context]({{< ref "content-blocks.md#context" >}}) for more details).


## Supported Arguments

The arguments provided in the block are either generic arguments or plugin-specific input parameters.


### Generic Arguments

- `config`: (optional) a reference to a named config block defined outside the document. If provided, it takes precedence over the default configuration for the plugin. See [Plugin Configuration]({{< ref "configs.md#plugin-configuration" >}}) for the details.


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
