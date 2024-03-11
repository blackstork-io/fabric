---
title: stixview 
plugin:
  name: blackstork/stixview
  description: ""
  tags: []
  version: "v0.4.0"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/stixview/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/stixview" "stixview" "v0.4.0" "stixview" "content provider" >}}

## Installation

To use `stixview` content provider, you must install the plugin `blackstork/stixview`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/stixview" = ">= v0.4.0"
  }
}
```

Note the version constraint set for the plugin.


#### Configuration

The content provider doesn't support any configuration parameters.

#### Usage

The content provider supports the following execution parameters:

```hcl
content stixview {
    caption = <string>  # optional
    gist_id = <string>  # optional
    height = <number>  # optional
    show_footer = <bool>  # optional
    show_idrefs = <bool>  # optional
    show_labels = <bool>  # optional
    show_marking_nodes = <bool>  # optional
    show_sidebar = <bool>  # optional
    show_tlp_as_tags = <bool>  # optional
    stix_url = <string>  # optional
    width = <number>  # optional
}
```

