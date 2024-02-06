---
title: Data Blocks
type: docs
weight: 1
---

# Data Blocks

The `data` block defines a call to a data plugin, fetching external data for content rendering. The order of `data` block definitions in the document is inconsequential.

```hcl
data <plugin-name> "<result-name>" {
  ...
}

document "foobar" {

  data <plugin-name> "<result-name>" {
    ...
  }

}
```

If the `data` block is declared at the root level of the configuration file, both names—the plugin name and the result name—must be provided. This pair acts as a unique identifier within the codebase, crucial for referencing this block. `data` blocks defined outside the document are not executed independently but must be referenced inside the document template.

However, if the `data` block is defined within the document, it must be at the root level. The plugin name and result name are required, forming a unique pair within the document's scope, serving as an identifier.

The arguments set in the block include input parameters for the data plugin, along with the plugin configuration. The data returned by the plugin is stored in the global context map under the path `data.<plugin-name>.<result-name>`.

## Supported Arguments

- `config`: _(optional)_ a reference to a named config block defined at the root level. If provided, it takes precedence over the default config for the plugin.

Other supported arguments may be defined by the plugin.

## Supported Nested Blocks

- `meta`
- `config`: _(optional)_ an inline config block for the plugin. If provided, it takes precedence over the default configuration for the plugin.

Other nested blocks are not supported.

## References

For more information related to references, see [here](../reference.md#data-block-references).