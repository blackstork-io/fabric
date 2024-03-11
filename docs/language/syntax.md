---
title: Syntax
type: docs
weight: 10
---

# Syntax

This page describes the native syntax of the Fabric Configuration Language (FCL). Leveraging the foundation laid by the [HashiCorp Configuration Language](https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md) (HCL), FCL aligns itself with a syntax favored by numerous applications for its simplicity, readability, and clarity.

The syntax of the Fabric language revolves around two fundamental components: arguments and blocks. These elements constitute the building blocks for crafting configurations within the Fabric Configuration Language.

## Arguments

The arguments play a crucial role in assigning values to names within a block. An example of using arguments is as follows:

```hcl
... {
    query_string = "kibana.alert.severity:critical"
}
```

The argument name, `query_string` in the snippet above, is allowed to contain letters, digits, underscores (`_`), and hyphens (`-`). However, the first character of an identifier must not be a digit.

## Blocks

In the Fabric Configuration Language (FCL), a block serves as a versatile container defining configurations, data requirements, or content structures. An example is provided below:

```hcl
document "alerts_overview" {

    content text {
        text = "Static text"
    }
    ...

}
```

A block is characterized by a type (`document` and `content` in the example above), dictating the number of labels permissible in a block signature. Additionally, a block can either bear a name (for example, "`alerts_overview`") or remain anonymous, as seen in the case of content text in the provided snippet. This flexibility in block composition contributes to the expressive and modular nature of FCL configurations.

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

## Next steps

See [Configuration]({{< ref "configs.md" >}}) documentation to learn how to configure Fabric and Fabric plugins.
