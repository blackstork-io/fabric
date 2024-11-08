---
title: blackstork/iris
weight: 20
plugin:
  name: blackstork/iris
  description: "The `iris` plugin for Iris Incident Response platform."
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/iris/"
type: docs
hideInMenu: true
---

{{< plugin-header "blackstork/iris" "iris" "v0.4.2" >}}

## Description
The `iris` plugin for Iris Incident Response platform.

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/iris" = ">= v0.4.2"
  }
}
```


## Data sources

{{< plugin-resources "iris" "data-source" >}}
