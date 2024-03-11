---
title: blackstork/stixview
weight: 20
plugin:
  name: blackstork/stixview
  description: ""
  tags: []
  version: "v0.4.0"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/stixview/"
type: docs
---

{{< plugin-header "blackstork/stixview" "stixview" "v0.4.0" >}}

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/stixview" = ">= v0.4.0"
  }
}
```



## Content providers

{{< plugin-resources "stixview" "content-provider" >}}