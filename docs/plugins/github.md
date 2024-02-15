---
title: blackstork/github
weight: 20
type: docs
---

# `blackstork/github` plugin

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/github" = "=> v0.0.0-dev"
  }
}
```

## Data sources

The plugin has the following data sources available:

### `github_issues`

#### Configuration

The data source supports the following configuration parameters:

```hcl
config data github_issues {
    github_token = <string>  # required
}
```

#### Usage

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