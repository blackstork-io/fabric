---
title: local_file
plugin:
  name: blackstork/builtin
  description: "Publishes content to local file"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: publisher
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "local_file" "publisher" >}}

The publisher is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The publisher doesn't support any configuration parameters.

#### Usage

The publisher supports the following execution parameters:

```hcl
publish local_file {
  # Path to the file
  #
  # Required string. For example:
  path = "dist/output.md"
}

```

