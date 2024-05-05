---
title: list
plugin:
  name: blackstork/builtin
  description: "Produces a list of items"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "list" "content provider" >}}

## Description
Produces a list of items

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration parameters.

#### Usage

The content provider supports the following execution parameters:

```hcl
content list {
  # Go template for the item of the list
  #
  # Required string.
  # For example:
  item_template = "[{{.Title}}]({{.URL}})"

  # Can be one of: "unordered", "ordered", "tasklist"
  #
  # Optional string.
  # Default value:
  format = "unordered"
}
```

