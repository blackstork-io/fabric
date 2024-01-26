---
title: Documents
type: docs
weight: 3 
---

# Documents

`document` block defines a document template with a specific name. The template name must be unique within the codebase, as it can be used as an identifier when referencing this block.

The document name is required. The document blocks must be defined on a root level of the configuration file and can not be inside other blocks.

```hcl
document <document-name> {
  ...
}
```

## Supported Arguments

- `title` – _(optional)_ a title of the document. It is a syntax sugar for a nested `content` block that renders a title. The title content block precedes any other nested `content` blocks or `sequence` blocks defined at the same level.

No other arguments are supported.

## Supported Nested Blocks

- `meta`
- `data`
- `content`
- `section`

