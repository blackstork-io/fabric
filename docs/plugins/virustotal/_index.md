---
title: blackstork/virustotal
weight: 20
plugin:
  name: blackstork/virustotal
  description: ""
  tags: []
  version: "v0.0.0-dev"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/virustotal/"
type: docs
---

{{< plugin-header "blackstork/virustotal" "virustotal" "v0.0.0-dev" >}}

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/virustotal" = ">= v0.0.0-dev"
  }
}
```

## Data sources

- [`virustotal_api_usage`]({{< relref "./data-sources/virustotal_api_usage" >}})
