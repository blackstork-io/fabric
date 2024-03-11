---
title: github_issues
plugin:
  name: blackstork/github
  description: ""
  tags: []
  version: "v0.4.0"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/github/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/github" "github" "v0.4.0" "github_issues" "data source" >}}

## Installation

To use `github_issues` data source, you must install the plugin `blackstork/github`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/github" = ">= v0.4.0"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data github_issues {
    github_token = <string>  # required
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data github_issues {
    assignee = <string>  # optional
    creator = <string>  # optional
    direction = <string>  # optional
    labels = <list of string>  # optional
    limit = <number>  # optional
    mentioned = <string>  # optional
    milestone = <string>  # optional
    repository = <string>  # required
    since = <string>  # optional
    sort = <string>  # optional
    state = <string>  # optional
}
```