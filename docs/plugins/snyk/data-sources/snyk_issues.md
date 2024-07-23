---
title: "`snyk_issues` data source"
plugin:
  name: blackstork/snyk
  description: "The `snyk_issues` data source fetches issues from Snyk"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/snyk/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/snyk" "snyk" "v0.4.2" "snyk_issues" "data source" >}}

## Description
The `snyk_issues` data source fetches issues from Snyk.

## Installation

To use `snyk_issues` data source, you must install the plugin `blackstork/snyk`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/snyk" = ">= v0.4.2"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration arguments:

```hcl
config data snyk_issues {
  # The Snyk API key
  #
  # Required string.
  # Must be non-empty
  # For example:
  api_key = "some string"
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data snyk_issues {
  # The group ID
  #
  # Optional string.
  # Default value:
  group_id = null

  # The organization ID
  #
  # Optional string.
  # Default value:
  org_id = null

  # The scan item ID
  #
  # Optional string.
  # Default value:
  scan_item_id = null

  # The scan item type
  #
  # Optional string.
  # Default value:
  scan_item_type = null

  # The issue type
  #
  # Optional string.
  # Default value:
  type = null

  # The updated before date
  #
  # Optional string.
  # Default value:
  updated_before = null

  # The updated after date
  #
  # Optional string.
  # Default value:
  updated_after = null

  # The created before date
  #
  # Optional string.
  # Default value:
  created_before = null

  # The created after date
  #
  # Optional string.
  # Default value:
  created_after = null

  # The effective severity level
  #
  # Optional list of string.
  # Default value:
  effective_severity_level = null

  # The status
  #
  # Optional list of string.
  # Default value:
  status = null

  # The ignored flag
  #
  # Optional bool.
  # Default value:
  ignored = null

  # The limit of issues to fetch
  #
  # Optional number.
  # Must be >= 0
  # Default value:
  limit = 0
}
```