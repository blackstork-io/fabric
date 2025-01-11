---
title: "`frontmatter` content provider"
plugin:
  name: blackstork/builtin
  description: "Produces the frontmatter"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.2" "frontmatter" "content provider" >}}

## Description
Produces the frontmatter.

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration arguments.

#### Usage

The content provider supports the following execution arguments:

```hcl
content frontmatter {
  # Format of the frontmatter.
  #
  # Optional string.
  # Must be one of: "yaml", "toml", "json"
  # Default value:
  format = "yaml"

  # Arbitrary key-value map to be put in the frontmatter.
  #
  # Required jq queriable.
  # Must be non-empty
  #
  # For example:
  content = {
    key = "arbitrary value"
    key2 = {
      "can be nested" = 42
    }
  }
}
```

