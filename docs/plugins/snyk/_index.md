---
title: blackstork/snyk
weight: 20
plugin:
  name: blackstork/snyk
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/snyk/"
type: docs
---

{{< plugin-header "blackstork/snyk" "snyk" "v0.4.1" >}}

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/snyk" = ">= v0.4.1"
  }
}
```


## Data sources

{{< plugin-resources "snyk" "data-source" >}}
