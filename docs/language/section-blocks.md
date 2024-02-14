---
title: Section Blocks
type: docs
weight: 70
---

# Section blocks

The blocks of type `section` are used for grouping and embedding `content` blocks. Section blocks can be referenced and nested within other section blocks.
Building a template with sections improves clarity and re-usability.

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

The section blocks respect the order of definition, same as the content blocks.

## Supported arguments

- `title`: (optional) represents the title of the content group. It's a syntactic sugar for a `content` block that renders a title. The title content block takes precedence over any other nested `content` blocks or `section` blocks defined at the same level.

## Supported nested blocks

- `meta`: (optional) a block containing metadata for the block.
- `content`: see [Content Blocks]({{< ref content-blocks.md >}}) for the details.
- `section`: a section block of can placed inside another `section` block.

## References

See [References]({{< ref references.md >}}) for the details about referencing section blocks.
