---
title: "`local_file` publisher"
plugin:
  name: blackstork/builtin
  description: "Publishes content to local file"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: publisher
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.2" "local_file" "publisher" >}}

The publisher is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

#### Formats

The publisher supports the following document formats:

- `md`
- `pdf`
- `html`

To set the output format, specify it inside `publish` block with `format` argument.


#### Configuration

The publisher doesn't support any configuration arguments.

#### Usage

The publisher supports the following execution arguments:

```hcl
# In addition to the arguments listed, `publish` block accepts `format` argument.

publish local_file {
  # Path to the file
  #
  # Required string.
  #
  # For example:
  path = "dist/output.md"
}

```

