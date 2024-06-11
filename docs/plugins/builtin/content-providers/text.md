---
title: "`text` content provider"
plugin:
  name: blackstork/builtin
  description: "Renders text"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "text" "content provider" >}}

## Description
Renders text

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration arguments.

#### Usage

The content provider supports the following execution arguments:

```hcl
content text {
  # A string to render. Can use go template syntax.
  #
  # Required string.
  # For example:
  value = "Hello world!"
}
```

