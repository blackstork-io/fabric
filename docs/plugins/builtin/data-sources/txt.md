---
title: txt
plugin:
  name: blackstork/builtin
  description: ""
  tags: []
  version: "v0.0.0-dev"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.0.0-dev" "txt" "data source" >}}

The data source is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

## Configuration

The data source doesn't support configuration.

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data txt {
    path = <string>  # required
}
```