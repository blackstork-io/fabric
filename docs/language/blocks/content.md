---
title: Content Blocks
type: docs
weight: 2
---

# Content Blocks

The `content` block defines a call to a content plugin responsible for producing a document segment. The sequence in which `content` blocks are defined is crucial as it determines the order in which the generated content will be placed in the document.

```hcl
# Root-level definition of the content block
content <plugin-name> "<block-name>" {
  ...
}

document "foobar" {

  # In-document definition of the content block
  content <plugin-name> "<block-name>" {
    ...
  }

  content <plugin-name> {
    ...
  }

}
```

If the `content` block is declared at the root level of the configuration file, both names—the plugin name and the block name—must be provided. This pair serves as a unique identifier within the codebase, crucial for referencing this block. `content` blocks defined outside the document are not executed independently but must be referenced inside the document template.

However, if the `content` block is defined within the document, only the plugin's name is required, and the block name becomes optional.

The arguments set in the block encompass input parameters for the content plugin, along with the plugin configuration and the local context map. Content plugins are designed to return a Markdown text string.

The order in which `content` blocks are defined is preserved.

## Supported Arguments

- `config`: _(optional)_ a reference to a named config block defined at the root level. If provided, it takes precedence over the default config for the plugin.
- `query`: _(optional)_ a `jq` query applied to the global context map; the result is stored under the `query_result` path in the local context map (a local extended copy of the global context map).
- `render_if_no_query_result`: _(optional)_ determines if the content block should be rendered when `query_result` is empty. Defaults to `true` (see [#28](https://github.com/blackstork-io/fabric/issues/28)).
- `text_when_no_query_result`: _(optional)_ provides text to be rendered instead of the plugin-returned text. Only used when `render_if_no_query_result` is `true` (see [#28](https://github.com/blackstork-io/fabric/issues/28)).
- `query_input_required`: _(optional)_ an attribute with a boolean value, set to false by default (see [#29](https://github.com/blackstork-io/fabric/issues/29)).
- `query_input`: _(optional)_ an attribute with a string value, empty by default (see [#29](https://github.com/blackstork-io/fabric/issues/29)).

Other supported arguments may be defined by the plugin.

## Supported Nested Blocks

- `meta`: _(optional)_ a block containing metadata for the `content` block.
- `config`: an inline config block for the plugin. If provided, it takes precedence over the default configuration for the plugin.

Other nested blocks are not supported.

## References

For more information related to references, see [here](../reference.md#content-block-references).