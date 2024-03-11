---
title: blackstork/sqlite
weight: 20
plugin:
  name: blackstork/sqlite
  description: ""
  tags: []
  version: "v0.0.0-dev"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/sqlite/"
type: docs
---

{{< plugin-header "blackstork/sqlite" "sqlite" "v0.0.0-dev" >}}

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/sqlite" = ">= v0.0.0-dev"
  }
}
```

## Data sources

- [`sqlite`]({{< relref "./data-sources/sqlite" >}})
