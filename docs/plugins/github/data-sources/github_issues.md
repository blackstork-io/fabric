---
title: "`github_issues` data source"
plugin:
  name: blackstork/github
  description: ""
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/github/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/github" "github" "v0.4.2" "github_issues" "data source" >}}

## Installation

To use `github_issues` data source, you must install the plugin `blackstork/github`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/github" = ">= v0.4.2"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration arguments:

```hcl
config data github_issues {
  # The GitHub token to use for authentication
  #
  # Required string.
  # Must be non-empty
  #
  # For example:
  github_token = "some string"
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data github_issues {
  # The repository to list issues from, in the format of owner/name
  #
  # Required string.
  # Must be non-empty
  #
  # For example:
  repository = "blackstork-io/fabric"

  # Filter issues by milestone. Possible values are:
  # * a milestone number
  # * "none" for issues with no milestone
  # * "*" for issues with any milestone
  # * "" (empty string) performs no filtering
  #
  # Optional string.
  # Default value:
  milestone = ""

  # Filter issues based on their state
  #
  # Optional string.
  # Must be one of: "open", "closed", "all"
  # Must be non-empty
  # Default value:
  state = "open"

  # Filter issues based on their assignee. Possible values are:
  # * a user name
  # * "none" for issues that are not assigned
  # * "*" for issues with any assigned user
  # * "" (empty string) performs no filtering.
  #
  # Optional string.
  # Default value:
  assignee = ""

  # Filter issues based on their creator. Possible values are:
  # * a user name
  # * "" (empty string) performs no filtering.
  #
  # Optional string.
  # Default value:
  creator = ""

  # Filter issues to once where this username is mentioned. Possible values are:
  # * a user name
  # * "" (empty string) performs no filtering.
  #
  # Optional string.
  # Default value:
  mentioned = ""

  # Filter issues based on their labels.
  #
  # Optional list of string.
  # Default value:
  labels = null

  # Specifies how to sort issues.
  #
  # Optional string.
  # Must be one of: "created", "updated", "comments"
  # Must be non-empty
  # Default value:
  sort = "created"

  # Specifies the direction in which to sort issues.
  #
  # Optional string.
  # Must be one of: "asc", "desc"
  # Must be non-empty
  # Default value:
  direction = "desc"

  # Only show results that were last updated after the given time.
  # This is a timestamp in ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ.
  #
  # Optional string.
  # Must be non-empty
  # Default value:
  since = null

  # Limit the number of issues to return. -1 means no limit.
  #
  # Optional integer.
  # Must be >= -1
  # Default value:
  limit = -1
}
```