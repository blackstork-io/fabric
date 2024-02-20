---
title: blackstork/stixview
weight: 20
type: docs
---

# `blackstork/stixview` plugin

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/stixview" = "=> v0.0.0-dev"
  }
}
```



## Content providers
The plugin has the following content providers available:

### `stixview`

#### Configuration

The content provider doesn't support configuration.

#### Usage

The content source supports the following parameters in the content blocks:

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
