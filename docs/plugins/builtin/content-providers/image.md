---
title: image
plugin:
  name: blackstork/builtin
  description: "Inserts an image"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "image" "content provider" >}}

## Description
Inserts an image

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration parameters.

#### Usage

The content provider supports the following execution parameters:

```hcl
content image {
  # Required string. For example:
  src = "https://example.com/img.png"

  # For example:
  # alt = "Text description of the image"
  #
  # Optional string. Default value:
  alt = null
}
```

