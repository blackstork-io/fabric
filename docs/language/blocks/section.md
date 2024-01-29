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

If the label `ref` is used instead of `<section-name>`, the block references another `section` block defined on a root level. The name of the referer block is optional if the block is defined within the document. If the referer block (with label `ref`) is defined on a root level of the config file, the name is required.

```hcl

section "foo" {
  ...
}

document "overview" {

  section ref {
    base = section.foo
    ...
  }

}
```

If `title` argument is provided in the referer block, it takes precedence over the title defined in the referent block.

Every referer block must have a `base` attribute set, pointing to a block defined on a root level in the config file.

Referer blocks can not contain nested blocks.



