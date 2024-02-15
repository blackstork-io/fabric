---
title: Documents
type: docs
weight: 40
---

# Documents

Document blocks are the most important element of the Fabric configuration. `document` block represents a template and includes the data and content blocks that define the document.

```hcl
document "<document-name>" {
  ...
}
```

A block type `document` and a document name are used as an unique identifier for the document template within the codebase.
The document blocks must be on a root level of the configuration file and can not be inside other blocks.

The `document` block is a structure that groups the data definitions, the section, and the content blocks together.

## Supported arguments

- `title`: (optional) a title of the document. It's a syntax sugar for a nested `content` block that renders a title. The title content block precedes any other nested `content` blocks or `sequence` blocks defined at the same level.

## Supported nested blocks

- `meta`: (optional) a block containing metadata for the block.
- `data`: see [Data Blocks]({{< ref data-blocks.md >}}) for the details.
- `content`: see [Content Blocks]({{< ref content-blocks.md >}}) for the details.
- `section`: see [Section Blocks]({{< ref section-blocks.md >}}) for the details.
