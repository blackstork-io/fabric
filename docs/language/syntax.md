---
title: Syntax
type: docs
weight: 1
---

# Syntax

The Fabric language syntax, similar to [Terraform language syntax](https://developer.hashicorp.com/terraform/language/syntax/configuration), is built around arguments and blocks.

## Arguments

The block where the argument appears determines what arguments are supported and what value types are valid.
**TBD**: (HCL expressions and functions support)


## Blocks

You can expect to find a limited set of block types in a Fabric config file: `document`, `data`, `content`, `section`, `meta`, `fabric`, and `config`.

The block type defines the labels supported. `data`, `content`, and `config` blocks are plugin-specific and require a plugin name to be present as a label.

