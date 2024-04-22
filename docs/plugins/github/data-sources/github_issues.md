---
title: github_issues
plugin:
  name: blackstork/github
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/github/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/github" "github" "v0.4.1" "github_issues" "data source" >}}

## Installation

To use `github_issues` data source, you must install the plugin `blackstork/github`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/github" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data github_issues {
  # Required string. For example:
  github_token = "some string"
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data github_issues {
  # Required string. For example:
  repository = "some string"

  # Optional string. Default value:
  milestone = null

  # Optional string. Default value:
  state = null

  # Optional string. Default value:
  assignee = null

  # Optional string. Default value:
  creator = null

  # Optional string. Default value:
  mentioned = null

  # Optional list of string. Default value:
  labels = null

  # Optional string. Default value:
  sort = null

  # Optional string. Default value:
  direction = null

  # Optional string. Default value:
  since = null

  # Optional number. Default value:
  limit = null
}
```