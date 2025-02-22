---
title: "`terraform_state_local` data source"
plugin:
  name: blackstork/terraform
  description: ""
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/terraform/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/terraform" "terraform" "v0.4.2" "terraform_state_local" "data source" >}}

## Installation

To use `terraform_state_local` data source, you must install the plugin `blackstork/terraform`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/terraform" = ">= v0.4.2"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source doesn't support any configuration arguments.

## Usage

The data source supports the following execution arguments:

```hcl
data terraform_state_local {
  # Required string.
  #
  # For example:
  path = "some string"
}
```