---
title: Content Blocks
type: docs
weight: 60
---

# Content Blocks

Blocks of type `content` define document segments: text paragraphs, tables, graphs, lists, etc. The order in which `content` blocks are defined determines the order of placement of the generated content in the document.

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

If a `content` block is defined at the root level of the configuration file, outside of the `document`, both names – the plugin name and the block name – must be provided. A combination of block type `content`, plugin name and block name serves as a unique identifier of a block within the codebase.

If a `content` block is defined within the document, only plugin name is required, while block name is optional.

A content block is implemented by a content plugin specified by `<plugin-name>` in a block signature. See [Plugins]({{< ref "plugins.md" >}}) for the details on the plugins supported by Fabric.

## Supported Arguments

The arguments provided in the block are either generic arguments or plugin-specific input parameters.

### Generic Arguments

- `config`: (optional) a reference to a named config block defined outside the document. If provided, it takes precedence over the default configuration for the plugin. See [Plugin Configuration]({{< ref "configs.md#plugin-configuration" >}}) for the details.
- `query`: (optional) a [JQ](https://jqlang.github.io/jq/manual/) query to be applied to the context object. The results of the query are stored under `query_result` field in the context. See [Context](#context) object for the details.
- `render_if_no_query_result`: (optional) a boolean flag that determines if the content block should be rendered when `query` returned no data. Defaults to `true`.
- `text_when_no_query_result`: (optional) provides the text to be rendered instead of calling the plugin when `render_if_no_query_result` is `true`.
- `query_input_required`: (optional) a boolean flag that specifies if `query_input` must be explicitely provided when the content block is referenced. `false` by default. See [Query Input Requirement]({{< ref "references.md#query-input-requirement" >}}) for the details.
- `query_input`: (optional) a JQ query to be applied to the context object. The results of the query are stored under `query_input` field in the context. See [Query Input Requirement]({{< ref "references.md#query-input-requirement" >}}) for the details.

### Plugin-specific Arguments

Plugin-specific arguments are defined by a plugin specification. See [Plugins]({{< ref "plugins.md" >}}) for the details on the supported arguments per plugin.

## Supported Nested Blocks

- `meta`: (optional) a block containing metadata for the block.
- `config`: (optional) an inline configuration for the block. If provided, it takes precedence over the `config` argument and default configuration for the plugin.

## Context

When Fabric renders a content block, a corresponding content plugin is called. Along with [plugin configuration]({{< ref "configs/#plugin-configuration" >}}) and input parameters, a plugin receives the context object with all data available.

The context object is a JSON dictionary with pre-defined root-level fields:

- `data` contains a map of all resolved data definitions for the document. The JSON path to a specific data point follows the data block signature: `data.<plugin-name>.<block-name>`.
- `query_result` contains a result of the execution of a JQ query provided in `query` attribute, as requested by `query_input_required` attribute. This is mostly used in `ref` blocks to explicitely provide data to the query in `query` attribute.

## References

See [References]({{< ref references.md >}}) for the details about referencing content blocks.

## Example

```hcl

config content openai "test_account" {
  api_key = "openai-api-key"
}

document "test-document" {

   data inline "foo" {
     d = {
       items = ["a", "b", "c"]
     }
   }

   content text {
     # Query contains a JQ query executed against the context
     query = ".data.inline.foo.d.items | length"

     # The result of the query is stored in the `query_result` field in the context.
     # The context is available for the templating engine inside the `content.text` plugin.
     text = "There are {{ .query_result }} items"
   }

   content openai {
     config = config.content.openai.test_account
     query = ".data.inline.foo.d"

     prompt = "Describe the items provided in the list"
   }
}
```
