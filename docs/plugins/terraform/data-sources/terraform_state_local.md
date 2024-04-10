---
title: terraform_state_local
plugin:
  name: blackstork/terraform
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/terraform/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/terraform" "terraform" "v0.4.1" "terraform_state_local" "data source" >}}

## Installation

To use `terraform_state_local` data source, you must install the plugin `blackstork/terraform`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/terraform" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source doesn't support configuration.

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data "terraform_state_local" {
  # Required. For example:
  path = "some string"
}

```