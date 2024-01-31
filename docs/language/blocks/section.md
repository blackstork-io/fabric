---
title: Section blocks
type: docs
weight: 3
---


# Section blocks ([#16](https://github.com/blackstork-io/fabric/issues/16))

`section` blocks are used for grouping `content` blocks for easier reusability and referencing. `section` blocks can contain other `section` blocks inside.

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

A section name must be provided if the `section` block is defined on a root level of the configuration file. The section name must be unique within the codebase, as it will be used as an identifier when referencing this block. The section blocks defined outside the document are not executed independently but must be referenced inside the document template.

If the `section` block is defined inside the document template, a section name is optional.


## Supported Arguments

- `title` – _(optional)_ a title of the content group. It is a syntax sugar for a nested `content` block that renders a title. The title content block precedes any other nested `content` blocks or `sequence` blocks defined at the same level.


No other arguments are supported.


## Supported Nested Blocks

- `meta`
- `content`
- `section`


## References ([#9](https://github.com/blackstork-io/fabric/issues/9))

To see information related to refreces see [here](../refrence.md#section-block-references).
