---
title: blackstork/postgresql
weight: 20
plugin:
  name: blackstork/postgresql
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/postgresql/"
type: docs
---

{{< plugin-header "blackstork/postgresql" "postgresql" "v0.4.1" >}}

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/postgresql" = ">= v0.4.1"
  }
}
```


## Data sources

{{< plugin-resources "postgresql" "data-source" >}}
