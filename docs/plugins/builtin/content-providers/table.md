---
title: table
plugin:
  name: blackstork/builtin
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "table" "content provider" >}}

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration parameters.

#### Usage

The content provider supports the following execution parameters:

```hcl
content "table" {
  # Required. For example:
  columns = [{
    header = "some string"
    value  = "some string"
    }, {
    header = "some string"
    value  = "some string"
  }]
}

```

