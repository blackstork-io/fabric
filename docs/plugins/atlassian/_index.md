---
title: blackstork/atlassian
weight: 20
plugin:
  name: blackstork/atlassian
  description: "The `atlassian` plugin for Atlassian Cloud."
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/atlassian/"
type: docs
hideInMenu: true
---

{{< plugin-header "blackstork/atlassian" "atlassian" "v0.4.2" >}}

## Description
The `atlassian` plugin for Atlassian Cloud.

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/atlassian" = ">= v0.4.2"
  }
}
```


## Data sources

{{< plugin-resources "atlassian" "data-source" >}}
