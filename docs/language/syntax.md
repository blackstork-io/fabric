---
title: Syntax
description: Explore the syntax of Fabric Configuration Language (FCL). Built upon the foundation of HCL, FCL offers simplicity, readability, and clarity. Learn about arguments and blocks, the fundamental components for crafting configurations with FCL. Dive into examples and understand the expressive and modular nature of FCL configurations.
type: docs
weight: 10
---

# Syntax

This page describes the native syntax of the Fabric Configuration Language (FCL). Building on the
foundation laid by the [HashiCorp Configuration Language](https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md) (HCL),
FCL defines a simple, readable, and clear syntax for document templates.

The syntax of FCL has two fundamental components: **arguments** and **blocks**. These components
constitute the building blocks for crafting configurations within the Fabric Configuration Language.

## Arguments

The arguments are used for assigning values to names within a block. An example of using arguments is as follows:

```hcl
content text {
  value = "An example of the text value"
}
```

The argument name, `value` in the snippet above, can contain letters, digits, underscores
(`_`), and hyphens (`-`). However, the first character of an identifier must not be a digit.

## Blocks

In the Fabric Configuration Language (FCL), a block serves as a container that defines
configurations, data requirements, content structures, variables, and publishing destinations.

For example:

```hcl
document "test_document" {

  data elasticsearch "events" {
    index = "events"
  }

  content text {
    value = "My custom static text"
  }

  vars {
    foo = "xyz"
  }

}
```

Each block has a type (for example, `document` and `content`) that defines the labels allowed in a
block signature. A block can have a name (for example, "`events`") or be anonymous, depending on the
type of the block and the position of the block in the code. This flexibility contributes to the
expressive and modular nature of FCL configurations.

The blocks types can be divided into the following categories based on their purpose:

- [Configuration]({{< ref "configs.md" >}}): `fabric`, `config`, and `meta` blocks
- [Documents]({{< ref "documents.md" >}}): `document` block
- [Data definitions]({{< ref "data-blocks.md" >}}): `data` block
- [Content definitions]({{< ref "content-blocks.md" >}}): `content` block
- [Content structure]({{< ref "section-blocks.md" >}}): `section` block
- [Publishing / Delivery]({{< ref "publish-blocks.md" >}}): `publish` block
- [Referenced blocks / Code reuse]({{< ref "references.md" >}}): reusing blocks by referencing

## Comments

Fabric language supports three different flavours of comments:

- `#` begins a single-line comment, ending at the end of the line.
- `//` is an alternative to `#` and also defines a single-line comment
- `/*` and `*/` are start and end delimiters for a comment that might span over multiple lines.

It's recommend to use `#` single-line comment style usually. Future Fabric code formatting tools
will prioritise `#` comments as idiomatic.

## Character encoding

Fabric configuration files must be UTF-8 encoded. Fabric supports non-ASCII characters in comments, and string values.

## Next steps

See [Configuration]({{< ref "configs.md" >}}) documentation to learn how to configure Fabric and Fabric plugins.
