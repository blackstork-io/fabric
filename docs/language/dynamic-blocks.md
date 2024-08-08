---
title: Dynamic Blocks
type: docs
weight: 80 
draft: true
---

# Dynamic blocks

Fabric templates define the structure of a document using `section` and `content` blocks. While a static document structure often suffices, there are situations where the document structure must depend on input data. `dynamic` blocks allow template authors to adjust the template structure based on the data.

Use `dynamic` blocks to enable the dynamic generation of `document`, `section`, and `content` blocks.

## Definition

A `dynamic` block extends the signature of the original block and adds its own arguments ([see
below](#dynamic-block-arguments)) to the original block's arguments and nested sub-blocks:

```hcl
dynamic <block-signature> {

  # Dynamic block arguments
  # ...

  # Original block arguments and nested blocks
  # ...
}
```

For example, this `dynamic` block, when placed in the document template, will produce 3 `content`
blocks:

```hcl
dynamic content text "foobarbaz" {

  # Dynamic block arguments
  dynamic_items = ["foo", "bar", "baz"]

  # The arguments below belong to a content block and will be evaluated after the execution
  # of the dynamic block
  var {
    x = 1
    item_upper = query_jq(".vars.dynamic_item | ascii_upcase")
  }

  value = <<-EOT
    Content block {{ .vars.dynamic_index }}:
    item={{ .vars.dynamic_item }} upper={{ .vars.item_upper }}
  EOT
}
```

The rendered output for it would be:

```text
Content block 0
item=foo upper=FOO

Content block 1
item=bar upper=BAR

Content block 2
item=baz upper=BAZ
```

## Arguments

- `dynamic_items`: (optional) a collection of items to iterate over. This can be a static value or a
  `query_jq()` function call.
- `dynamic_condition`: (optional) a boolean value. This can be a static value or a `query_jq()`
  function call.

Either `dynamic_items` value or `dynamic_condition` value must be provided.

- if `dynamic_items` is set, the `dynamic` block will produce one block for each item in the
  `dynamic_items` collection.
- if `dynamic_condition` is set, the `dynamic` block will produce one block if the condition is
  `true` and none if it is `false`.

## Extended context

The context for each block produced by a `dynamic` block with the `dynamic_items` argument includes
additional variables:

- `.vars.dynamic_item` — the value of the current iteration item
- `.vars.dynamic_index` — the index of the current iteration item

## Dynamic reference blocks

To create flexible and reusable templates, `dynamic` blocks can be combined with the reference blocks.

For example:

```hcl
dynamic section ref {
  dynamic_items = query_jq(".vars.defined_items")
  base = section.external_section
}
```

This dynamic block creates a `section ref` block for each item in the `dynamic_items` collection.

Dynamic blocks in Fabric templates offer the flexibility to create complex, data-driven documents. By using these blocks, template authors can generate dynamic content that adapts to the shape of input data.
