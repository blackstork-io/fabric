---
title: blackstork/microsoft
weight: 20
plugin:
  name: blackstork/microsoft
  description: "The `microsoft` plugin for Microsoft services."
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/microsoft/"
type: docs
hideInMenu: true
---

{{< plugin-header "blackstork/microsoft" "microsoft" "v0.4.2" >}}

## Description
The `microsoft` plugin for Microsoft services.

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/microsoft" = ">= v0.4.2"
  }
}
```


## Data sources

{{< plugin-resources "microsoft" "data-source" >}}
