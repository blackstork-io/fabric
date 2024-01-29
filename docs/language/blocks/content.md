---
title: Content blocks
type: docs
weight: 2
---

# Content blocks

`content` block defines a call to a content plugin that produces a document segment. The order in which the `content` blocks are defined is important – it is the order in which the generated content will be placed in the document.

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

If `content` block is defined on a root level of the configuration file, both names - the name of the plugin and the block name - must be provided. The plugin name / block name pair must be unique within the codebase, as it will be used as an identifier when referencing this block. The content blocks defined outside the document are not executed independently but must be referenced inside the document template.

If the `content` block is defined in the document, only the plugin's name is required, while the block name is optional.

The arguments set in the block are the input parameters for the content plugin, together with the plugin configuration and the local context map. Content plugins are expected to return a Markdown text string.

The order in which `content` blocks are defined is preserved.

## Supported Arguments

- `config` – _(optional)_ a reference to a named config block defined on a root level. If provided, it takes precedence over the default config for the plugin.
- `query` – _(optional)_ a `jq` query that will be applied to the global context map, and the result will be stored under the `query_result` path in the local context map (a local extended copy of the global context map).
- `render_if_no_query_result` – _(optional)_ if content block should be rendered while `query_result` is empty. `true` by default (see [#28](https://github.com/blackstork-io/fabric/issues/28))
- `text_when_no_query_result` – _(optional)_ provides a text to be rendered instead of the plugin-returned text. Only used when `render_if_no_query_result` is `true` (see [#28](https://github.com/blackstork-io/fabric/issues/28))
- `query_input_required` – _(optional)_ an attribute with a boolean value, set to false by default (see [#29](https://github.com/blackstork-io/fabric/issues/29))
- `query_input` – (optional) an attribute with a string value, empty by default. (see [#29](https://github.com/blackstork-io/fabric/issues/29))

The plugin might define other supported arguments.


## Supported Nested Blocks

- `meta` – _(optional)_ a block that contains metadata for the content block
- `config` – an inline config block for the plugin. If provided, it takes precedence over the default configuration for the plugin.

Other nested blocks are not supported.


## References

If the label `ref` is used instead of `<plugin-name>`, the block references another content block defined on a root level. The name of the referer block is optional if the block is defined within the document. If the referer block (with label `ref`) is defined on a root level of the config file, the name is required.

```hcl
content openai "foo" {
  ...
}

document "overview" {

  content ref {
    base = content.openai.foo
    ...
  }

}

```
Every referer data block must have a `base` attribute set, pointing to a data block defined on a root level in the config file.

Other supported arguments are defined by the plugin of the referent block (`openai` in the above example). If any arguments are provided in the referer block, they take precedence over the arguments set in the referent block.

Referer blocks can not contain nested blocks.
