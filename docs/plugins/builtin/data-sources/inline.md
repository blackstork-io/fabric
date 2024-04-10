---
title: inline
plugin:
  name: blackstork/builtin
  description: "Creates a queryable key-value map from the block's contents"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "inline" "data source" >}}

## Description
Creates a queryable key-value map from the block's contents

The data source is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

## Configuration

The data source doesn't support configuration.

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data "inline" {
  # Arbitrary structure of (possibly nested) blocks and attributes.
  # For example:
  #   key1 = "value1"
  #   nested {
  #     blocks {
  #       key2 = 42
  #     }
  #   }

}
```