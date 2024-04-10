---
title: stixview
plugin:
  name: blackstork/stixview
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/stixview/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/stixview" "stixview" "v0.4.1" "stixview" "content provider" >}}

## Installation

To use `stixview` content provider, you must install the plugin `blackstork/stixview`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/stixview" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.


#### Configuration

The content provider doesn't support any configuration parameters.

#### Usage

The content provider supports the following execution parameters:

```hcl
content "stixview" {
  # Optional. Default value:
  gist_id = null

  # Optional. Default value:
  stix_url = null

  # Optional. Default value:
  caption = null

  # Optional. Default value:
  show_footer = null

  # Optional. Default value:
  show_sidebar = null

  # Optional. Default value:
  show_tlp_as_tags = null

  # Optional. Default value:
  show_marking_nodes = null

  # Optional. Default value:
  show_labels = null

  # Optional. Default value:
  show_idrefs = null

  # Optional. Default value:
  width = null

  # Optional. Default value:
  height = null
}
```

