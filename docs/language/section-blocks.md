---
title: Section Blocks
type: docs
weight: 70
---

# Section Blocks

The blocks of type `section` are used for grouping and nesting `content` blocks for better clarity and reusability.

```hcl
section "<section-name>" {
  ...
}

document "foobar" {

  section {
    ...
  }

  section "<section-name>" {

    section {
      ...
    }

    ...
  }
}
```

If a `section` block is defined at the root level of the configuration file, outside of the `document`, the section name is required. A combination of a block type (`section`) and a section name serves as a unique identifier of a block within the codebase.

If a `section` block is defined within the document template, the section name is optional.


## Supported Arguments

- `title`: (optional) represents the title of the content group. It is a syntactic sugar for a `content` block that renders a title. The title content block takes precedence over any other nested `content` blocks or `section` blocks defined at the same level.


## Supported Nested Blocks

- `meta`: (optional) a block containing metadata for the block.
- `content`: see [Content Blocks]({{< ref content-blocks.md >}}) for the details.
- `section`: a block of type `section` can be embedded into another `section` block.


## References

See [References]({{< ref references.md >}}) for the details about referencing section blocks.
