---
title: Section Blocks
description: Learn how to use use section blocks to build modular and reusable content structures in your templates.
type: docs
weight: 70
---

# Section blocks

`section` blocks group and provide scope for the document content.

Section blocks can be referenced, can contain metadata and can be nested within other section
blocks. Building a template with sections improves clarity and re-usability.

```hcl
section "<section-name>" {
  ...
}

document "soc-activity-overview" {

  section "ExecSummary" {
    ...
  }

  section "KPIs" {

    section "SLAs" {
      ...
    }

    section "Coverage" {
      ...
    }

    ...
  }
}
```

When a `section` block is defined at the root level of the configuration file, outside of the `document` block, the section name is required. A combination of a block type (`section`) and a section name serves as a unique identifier for a block within the codebase.

If a `section` block is defined within the `document` block, the section name is optional.

Similarly to the `content` blocks, the `section` blocks are rendered in the order of definition.

## Supported arguments

- `title`: (optional) represents the title of the content group. It's a syntactic sugar for a
  `content` block that renders a title. The title content block takes precedence over any other
  nested `content` blocks or `section` blocks defined at the same level.
- `local_var`: (optional) a shortcut for specifying a local variable. See [Variables]({{< ref
  "context.md#variables" >}}) for the details.

## Supported nested blocks

- `meta`: (optional) a block containing metadata for the block. See [Metadata]({{< ref "configs.md#metadata" >}}) for details.
- `content`: see [Content Blocks]({{< ref content-blocks.md >}}) for the details.
- `section`: nested `section` blocks.

## References

See [References]({{< ref references.md >}}) for the details about referencing section blocks.
