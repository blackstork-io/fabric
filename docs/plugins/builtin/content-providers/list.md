---
title: "`list` content provider"
plugin:
  name: blackstork/builtin
  description: "Produces a list of items"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.2" "list" "content provider" >}}

## Description
Produces a list of items

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration arguments.

#### Usage

The content provider supports the following execution arguments:

```hcl
content list {
  # Go template for the item of the list
  #
  # Optional string.
  #
  # For example:
  # item_template = "[{{.Title}}]({{.URL}})"
  #
  # Default value:
  item_template = "{{.}}"

  # Optional string.
  # Must be one of: "unordered", "ordered", "tasklist"
  # Default value:
  format = "unordered"

  # List of items to render.
  #
  # Required list of jq queriable.
  # Must be non-empty
  #
  # For example:
  items = ["First item", "Second item", "Third item"]
}
```

