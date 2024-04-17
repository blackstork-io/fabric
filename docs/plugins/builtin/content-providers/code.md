---
title: code
plugin:
  name: blackstork/builtin
  description: "Formats text as code snippet"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "code" "content provider" >}}

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration parameters.

#### Usage

The content provider supports the following execution parameters:

```hcl
content code {
  # Required string. For example:
  value = "Text to be formatted as a code block"

  # Specifiy the language for syntax highlighting
  #
  # For example:
  # language = "python3"
  #
  # Optional string. Default value:
  language = ""
}
```

