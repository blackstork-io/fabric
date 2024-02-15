---
title: Syntax
type: docs
weight: 10
---

# Syntax

This page describes native syntax of the Fabric Configuration Language (FCL). FCL is based on [HCL](https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md) (HashiCorp Configuration Language) favored by many applications for simplicity, readability, and clarity.

Fabric language syntax has two core components: arguments and blocks.

## Arguments

The arguments are used for assigning a value to a name inside a block:

```hcl
... {
    query_string = "kibana.alert.severity:critical"
}
```

The argument name (`query_string` in the snippet above) can contain letters, digits, underscores (`_`), and hyphens (`-`). The first character of an identifier must not be a digit.

## Blocks

A block is a container that defines a configuration, a data requirement or a content structure.

```hcl
document "alerts_overview" {

    content text {
        text = "Static text"
    }
    ...

}
```

A block has a type (`document` and `content` in the example above) that defines how many labels can be used in a block signature. A block can have a name (`"alerts_overview"`) or be anonymous (as a `content text` in the snippet above).

Supported categories of blocks:

- [Configuration]({{< ref "configs.md" >}}): `fabric` and `config` blocks
- [Documents]({{< ref "documents.md" >}}): `document` block
- [Data definitions]({{< ref "data-blocks.md" >}}): `data` block
- [Content definitions]({{< ref "content-blocks.md" >}}): `content` block
- [Content structure]({{< ref "section-blocks.md" >}}): `section` block

## Comments

Fabric language supports three different flavours of comments:

- `#` begins a single-line comment, ending at the end of the line.
- `//` is an alternative to `#` and also defines a single-line comment
- `/*` and `*/` are start and end delimiters for a comment that might span over multiple lines.

It's recommend to use `#` single-line comment style usually. Future Fabric code formatting tools will prioritise `#` comments as idiomatic.

## Character encoding

Fabric configuration files must be UTF-8 encoded. Fabric allows non-ASCII characters in comments, and string values.
