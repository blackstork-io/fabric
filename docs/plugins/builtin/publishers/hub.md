---
title: "`hub` publisher"
plugin:
  name: blackstork/builtin
  description: "Publish documents to Blackstork Hub."
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: publisher
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.2" "hub" "publisher" >}}

The publisher is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

#### Formats

The publisher supports the following document formats:

- `unknown`

To set the output format, specify it inside `publish` block with `format` argument.


#### Configuration

The publisher supports the following configuration arguments:

```hcl
config publish hub {
  # API url.
  #
  # Required string.
  # Must be non-empty
  # For example:
  api_url = "some string"

  # API url.
  #
  # Required string.
  # Must be non-empty
  # For example:
  api_token = "some string"
}

```

#### Usage

The publisher supports the following execution arguments:

```hcl
# In addition to the arguments listed, `publish` block accepts `format` argument.

publish hub {
  # Hub Document title override. By default uses title configured in the document.
  #
  # Optional string.
  # Must be non-empty
  # Default value:
  title = null
}

```

