---
title: list 
plugin:
  name: blackstork/builtin
  description: ""
  tags: []
  version: "v0.4.0"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.0" "list" "content provider" >}}

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration parameters.

#### Usage

The content provider supports the following execution parameters:

```hcl
content list {
    format = <string>  # optional
    item_template = <string>  # required
}
```

