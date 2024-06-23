---
title: Evaluation Context
description: Learn about the evaluation context, the data structure that during rendering holds data available for the Fabric template blocks.
type: docs
weight: 45
---

# Evaluation context

The evaluation context keeps all available data during the template evaluation. This includes the
results of `data` block executions, the template metadata, and defined variables.

The context data can be queried with `query_jq()` function and accessed in Go templates (if
supported by the plugin).

## Keys

The context is a dictionary. To access the values in the context, use JSON paths in [JQ queries](https://jqlang.github.io/jq/manual/).

Using JQ-style JSON path syntax, the root keys are:

- `.data` — stores the results of the `data` blocks evaluation under
  `.data.<data-source-name>.<block-name>` keys
- `.document` — contains the data about the current document template. For example, `.document.meta`
  has values from the `meta` block defined on the root level of the template
- `.vars` — holds the values of variables defined in the current or parent blocks.

The context also can contain local properties:

- `.section` — holds the data about the current section, when called inside one. Similar to
  `.document.meta`, `.section.meta` contains the data from `meta` block defined inside the section.
- `.content` — holds the data about the current content block, when called inside one. Similar to
  `.document.meta` and `.section.meta`, `.content.meta` stores the data from `meta` block defined
  inside the content block.

## Variables

Define variables inside the template using the `vars` block:

```hcl
vars {
  foo = 1
  bar = {
    x = "a"
    y = env.CUSTOM_ENV_VAR
  }
  # The variables are evaluated in the order of definition
  baz = query_jq(".vars.foo")
}
```

The `vars` block can be defined inside `document`, `section`, `content`, dynamic and reference blocks.

The variable can hold static values (like `foo` in the example above) or the results of data
mutations (like `bar` that holds the result of JQ query executed in `query_jq()`).

After evaluation the variables become available in the context under the `.vars` root keyword.
For example, `.vars.foo` and `.vars.bar` JQ paths refer to the values of `foo` and `bar` variables.

### Local variable

To define a local variable, use `local_var` argument. It's a shortcut that defines a single local
variable named `local` in the current block.

For example:

```hcl
content text {
  local_var = "World"
  value = "Hello, {{ .vars.local }}!"
}
```

content block will render `Hello, World!`.

### Inheritance

Variables defined in a parent block (such as `document` or `section`) are available for use in
nested blocks. Note that nested blocks can redefine a variable in their context, overwriting the
parent variable's value.

For example:

```hcl
section {
  vars {
    foo = 11
    bar = 22
  }

  section {
    vars {
      foo = 33
      baz = 44
    }

    content text {
      # Renders: `Variable values: foo=33, bar=22, baz=44`
      value = "Variable values: foo={{ .vars.foo }}, bar={{ .vars.bar }}, baz={{ .vars.baz }}"
    }
  }
}
```

### Querying the context

To filter and mutate the data in the context, use [JQ queries](https://jqlang.github.io/jq/manual/).

`query_jq()` function accepts a string argument as a query, executes the query against the
evaluation context, and returns the results.

For example:

```hcl
section {
  vars {
    items = ["a", "b", "c"]
  }

  section {
    vars {
      items_count = query_jq(".vars.items | length")

      items_uppercase = query_jq(
        <<-EOT
          .vars.items | map(ascii_upcase) | join(":")
        EOT
      )
    }

    content text {
      # Renders: `Items count: 3; Uppercase items: A:B:C`
      value = "Items count: {{ .vars.items_count }}; Uppercase items: {{ .vars.items_uppercase }}"
    }
  }
}
```

## Next steps

See [Data Blocks]({{< ref "data-blocks.md" >}}) documentation to learn how to define data requirements in the templates.
