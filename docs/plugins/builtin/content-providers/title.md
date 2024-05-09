---
title: title
plugin:
  name: blackstork/builtin
  description: "Produces a title"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "title" "content provider" >}}

## Description
Produces a title.

The title size after calculations must be in an interval [0; 5] inclusive, where 0
corresponds to the largest size (`<h1>`) and 5 corresponds to (`<h6>`)

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration parameters.

#### Usage

The content provider supports the following execution parameters:

```hcl
content title {
  # Title content
  #
  # Required string.
  # For example:
  value = "Vulnerability Report"

  # Sets the absolute size of the title.
  # If `null` – absoulute title size is determined from the document structure
  #
  # Optional integer.
  # Default value:
  absolute_size = null

  # Adjusts the absolute size of the title.
  # The value (which may be negative) is added to the `absolute_size` to produce the final title size
  #
  # Optional integer.
  # Default value:
  relative_size = 0
}
```

