---
title: "`blockquote` content provider"
plugin:
  name: blackstork/builtin
  description: "Formats text as a block quote"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "blockquote" "content provider" >}}

## Description
Formats text as a block quote

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration arguments.

#### Usage

The content provider supports the following execution arguments:

```hcl
content blockquote {
  # Required string.
  # For example:
  value = "Text to be formatted as a quote"
}
```

