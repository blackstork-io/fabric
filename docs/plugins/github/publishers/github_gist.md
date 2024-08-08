---
title: "`github_gist` publisher"
plugin:
  name: blackstork/github
  description: "Publishes content to github gist"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/github/"
resource:
  type: publisher
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/github" "github" "v0.4.2" "github_gist" "publisher" >}}

## Installation

To use `github_gist` publisher, you must install the plugin `blackstork/github`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/github" = ">= v0.4.2"
  }
}
```

Note the version constraint set for the plugin.

#### Formats

The publisher supports the following document formats:

- `md`
- `html`

To set the output format, specify it inside `publish` block with `format` argument.


#### Configuration

The publisher supports the following configuration arguments:

```hcl
config publish github_gist {
  # Required string.
  # For example:
  github_token = "some string"
}

```

#### Usage

The publisher supports the following execution arguments:

```hcl
# In addition to the arguments listed, `publish` block accepts `format` argument.

publish github_gist {
  # Optional string.
  # Default value:
  description = null

  # Optional string.
  # Default value:
  filename = null

  # Optional bool.
  # Default value:
  make_public = false

  # Optional string.
  # Default value:
  gist_id = null
}

```

