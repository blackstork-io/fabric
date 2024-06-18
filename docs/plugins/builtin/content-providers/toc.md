---
title: "`toc` content provider"
plugin:
  name: blackstork/builtin
  description: "Produces table of contents"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "toc" "content provider" >}}

## Description
Produces table of contents.

Inspects the rendered document for headers of a certain size and creates a linked
table of contents

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration arguments.

#### Usage

The content provider supports the following execution arguments:

```hcl
content toc {
  # Largest header size which produces entries in the table of contents
  #
  # Optional integer.
  # Must be between 0 and 5 (inclusive)
  # Default value:
  start_level = 0

  # Smallest header size which produces entries in the table of contents
  #
  # Optional integer.
  # Must be between 0 and 5 (inclusive)
  # Default value:
  end_level = 2

  # Whether to use ordered list for the contents
  #
  # Optional bool.
  # Default value:
  ordered = false

  # Scope of the headers to evaluate.
  #   "document" – look for headers in the whole document
  #   "section" – look for headers only in the current section
  #   "auto" – behaves as "section" if the "toc" block is inside of a section; else – behaves as "document"
  #
  # Optional string.
  # Must be one of: "document", "section", "auto"
  # Default value:
  scope = "auto"
}
```

