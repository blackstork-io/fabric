---
title: "`image` content provider"
plugin:
  name: blackstork/builtin
  description: "Returns an image tag"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.2" "image" "content provider" >}}

## Description
Returns an image tag

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration arguments.

#### Usage

The content provider supports the following execution arguments:

```hcl
content image {
  # Required string.
  # Must have a length of at least 1
  # For example:
  src = "https://example.com/img.png"

  # Optional string.
  # For example:
  # alt = "Text description of the image"
  # 
  # Default value:
  alt = null
}
```

