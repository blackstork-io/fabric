---
title: "`misp_event_reports` publisher"
plugin:
  name: blackstork/misp
  description: "Publishes content to misp event reports"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/misp/"
resource:
  type: publisher
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/misp" "misp" "v0.4.2" "misp_event_reports" "publisher" >}}

## Installation

To use `misp_event_reports` publisher, you must install the plugin `blackstork/misp`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/misp" = ">= v0.4.2"
  }
}
```

Note the version constraint set for the plugin.

#### Formats

The publisher supports the following document formats:

- `md`

To set the output format, specify it inside `publish` block with `format` argument.


#### Configuration

The publisher supports the following configuration arguments:

```hcl
config publish misp_event_reports {
  # misp api key
  #
  # Required string.
  # Must be non-empty
  # For example:
  api_key = "some string"

  # misp base url
  #
  # Required string.
  # Must be non-empty
  # For example:
  base_url = "some string"

  # skip ssl verification
  #
  # Optional bool.
  # Default value:
  skip_ssl = false
}

```

#### Usage

The publisher supports the following execution arguments:

```hcl
# In addition to the arguments listed, `publish` block accepts `format` argument.

publish misp_event_reports {
  # Required string.
  # Must be non-empty
  # For example:
  event_id = "some string"

  # Required string.
  # Must be non-empty
  # For example:
  name = "some string"

  # Optional string.
  # Must be one of: "0", "1", "2", "3", "4", "5"
  # Default value:
  distribution = null

  # Optional string.
  # Default value:
  sharing_group_id = null
}

```

