---
title: Section Blocks
type: docs
weight: 3
---

# Section Blocks ([#16](https://github.com/blackstork-io/fabric/issues/16))

`section` blocks play a crucial role in grouping `content` blocks for enhanced reusability and referencing. Additionally, these versatile blocks can encapsulate other `section` blocks.

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

When declaring a `section` block at the root level of the configuration file, ensure a section name is provided. This name serves as a unique identifier within the codebase, vital for referencing this block. `section` blocks defined outside the document aren't executed independently but must be referenced inside the document template.

However, if a `section` block is defined within the document template, the section name becomes optional.

## Supported Arguments

- `title`: _(optional)_ represents the title of the content group. It acts as a syntactic sugar for a nested `content` block that renders a title. The title content block takes precedence over any other nested `content` blocks or `sequence` blocks defined at the same level.

No other arguments are supported.

## Supported Nested Blocks

- `meta`
- `content`
- `section`

## References ([#9](https://github.com/blackstork-io/fabric/issues/9))

For further details regarding references, visit [here](../reference.md#section-block-references).