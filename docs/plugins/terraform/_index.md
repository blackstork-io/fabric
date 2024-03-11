---
title: blackstork/terraform
weight: 20
plugin:
  name: blackstork/terraform
  description: ""
  tags: []
  version: "v0.4.0"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/terraform/"
type: docs
---

{{< plugin-header "blackstork/terraform" "terraform" "v0.4.0" >}}

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/terraform" = ">= v0.4.0"
  }
}
```


## Data sources

{{< plugin-resources "terraform" "data-source" >}}
