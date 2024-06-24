---
title: Documents
description: Learn how to use document blocks to define document templates, building a content structure with the content blocks and setting data requirements with the data blocks.
type: docs
weight: 40
---

# Documents

Document blocks are the most important element of the Fabric configuration. `document` block represents a template and includes the data, content and publishing definitions that describe the document and the workflow around it.

```hcl
document "<document-name>" {

  title = "<document title>"

  # ...
}
```

A block type `document` and a document name are an unique identifier for the document template within the codebase. The document blocks must defined be on a root level of the configuration file and can not be inside other blocks.

The `document` block is a structure that groups the data definitions, the sections, the content blocks, and the publishing instructions together.

## Supported arguments

- `title`: (optional) a title of the document. It's a syntax sugar for a nested `content` block that renders a title. During rendering, the title precedes any other nested `content` blocks or `section` blocks defined at the root level of the template.

## Supported nested blocks

- `meta`: see [Metadata]({{< ref "configs.md/#metadata" >}})
- `data`: see [Data Blocks]({{< ref data-blocks.md >}})
- `vars`: see [Variables]({{< ref "context.md/#variables" >}})
- `content`: see [Content Blocks]({{< ref content-blocks.md >}})
- `section`: see [Section Blocks]({{< ref section-blocks.md >}})
- `publish`: see [Publish Blocks]({{< ref publish-blocks.md >}})

## Next steps

See [Evaluation Context]({{< ref context.md >}}) documentation to learn how about the context that
holds all data available for the template.
