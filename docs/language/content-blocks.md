---
title: Content Blocks
type: docs
weight: 60
---

# Content blocks

`content` blocks define document segments: text paragraphs, tables, graphs, lists, etc. The order of the `content` blocks in the template determines the order of the generated content in the document.

```hcl
# Root-level definition of the content block
content <content-provider-name> "<block-name>" {
  ...
}

document "foobar" {

  # In-document definition of the content block
  content <content-provider-name> "<block-name>" {
    ...
  }

  content <content-provider-name> {
    ...
  }

}
```

If the block is defined at the root level of the file, outside of the `document` block, both names – the content provider name and the block name – are required. A combination of block type `content`, content provider name, and block name serves as a unique identifier of a block within the codebase.

If the content block is defined within the document template, only a content provider name is required and a block name is optional.

A content block is rendered by a specified content provider. See [Plugins]({{< ref "plugins.md" >}}) for the list of the content providers supported by Fabric.

## Supported arguments

The arguments provided in the block are either generic arguments or plugin-specific input parameters.

### Generic arguments

- `config`: (optional) a reference to a named configuration block for the data source or a content provider. If provided, it takes precedence over the default configuration. See [Content provider configuration]({{< ref "configs.md#content-provider-configuration" >}}) for the details.
- `query`: (optional) a [JQ](https://jqlang.github.io/jq/manual/) query to be executed against the context object. The results of the query will be placed under `query_result` field in the context. See [Context](#context) object for the details.
- `render_if_no_query_result`: (optional) a boolean flag that determines if the content block should be rendered when `query` returned no data. Defaults to `true`.
- `text_when_no_query_result`: (optional) provides the text to be rendered instead of calling the plugin when `render_if_no_query_result` is `true`.
- `query_input_required`: (optional) a boolean flag that specifies if `query_input` must be explicitly provided when the content block is referenced. `false` by default. See [Query Input Requirement]({{< ref "references.md#query-input-requirement" >}}) for the details.
- `query_input`: (optional) a JQ query to be applied to the context object. The results of the query are stored under `query_input` field in the context. See [Query Input Requirement]({{< ref "references.md#query-input-requirement" >}}) for the details.

### Content provider arguments

Content provider specific are defined by a plugin specification. See [Plugins]({{< ref "plugins.md" >}}) for the details on the supported arguments per plugin.

## Supported nested blocks

- `meta`: (optional) a block containing metadata for the block.
- `config`: (optional) an inline configuration for the block. If provided, it takes precedence over the `config` argument and default configuration for the plugin.

## Context

When Fabric renders a content block, a corresponding content plugin is called. Along with the [content provider configuration]({{< ref "configs/#content-provider-configuration" >}}) and input parameters, a plugin receives the context object with all data available.

The context object is a JSON dictionary with pre-defined root-level fields:

- `data` points to a map of all resolved data definitions for the document. The JSON path to a specific data point follows the data block signature: `data.<plugin-name>.<block-name>`.
- `query_result` points to a result of the execution of a JQ query provided in `query` attribute, as requested by `query_input_required` attribute. This is mostly used in `ref` blocks to explicitly pre-filter the data for the JQ query set in `query` attribute.

## References

See [References]({{< ref references.md >}}) for the details about referencing content blocks.

## Example

```hcl
config content openai_text "test_account" {
  api_key = "<OPENAI-KEY>"
}

document "test-doc" {

  data inline "foo" {
    items = ["aaa", "bbb", "ccc"]
  }

  content text {
    # Query contains a JQ query executed against the context
    query = ".data.inline.foo.items | length"

    # The result of the query is stored in the `query_result` field in the context.
    # The context is available for the templating engine inside the `content.text` plugin.
    text = "There are {{ .query_result }} items"
  }

  content openai_text {
    config = config.content.openai_text.test_account
    query = ".data.inline.foo"

    prompt = <<-EOT
       Write a short story, just a paragraph, about space exploration
       using the values from the provided items list as character names.
    EOT
  }
}
```

produces the output
```
There are 3 items

In the vast expanse of the universe, three brave astronauts, aaa, bbb, and ccc, embarked on a daring mission of space exploration. As they soared through the galaxies, their unwavering determination and unyielding teamwork propelled them towards uncharted territories, uncovering hidden wonders and pushing the boundaries of human understanding. Together, aaa, bbb, and ccc, etched their names in the stars as pioneers of a new era, forever inspiring generations to dream beyond the confines of Earth.
```
