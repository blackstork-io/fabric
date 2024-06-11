---
title: "`hackerone_reports` data source"
plugin:
  name: blackstork/hackerone
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/hackerone/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/hackerone" "hackerone" "v0.4.1" "hackerone_reports" "data source" >}}

## Installation

To use `hackerone_reports` data source, you must install the plugin `blackstork/hackerone`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/hackerone" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration arguments:

```hcl
config data hackerone_reports {
  # Required string.
  # For example:
  api_username = "some string"

  # Required string.
  # For example:
  api_token = "some string"
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data hackerone_reports {
  # Optional number.
  # Default value:
  size = null

  # Optional number.
  # Default value:
  page_number = null

  # Optional string.
  # Default value:
  sort = null

  # Optional list of string.
  # Default value:
  program = null

  # Optional list of number.
  # Default value:
  inbox_ids = null

  # Optional list of string.
  # Default value:
  reporter = null

  # Optional list of string.
  # Default value:
  assignee = null

  # Optional list of string.
  # Default value:
  state = null

  # Optional list of number.
  # Default value:
  id = null

  # Optional list of number.
  # Default value:
  weakness_id = null

  # Optional list of string.
  # Default value:
  severity = null

  # Optional bool.
  # Default value:
  hacker_published = null

  # Optional string.
  # Default value:
  created_at__gt = null

  # Optional string.
  # Default value:
  created_at__lt = null

  # Optional string.
  # Default value:
  submitted_at__gt = null

  # Optional string.
  # Default value:
  submitted_at__lt = null

  # Optional string.
  # Default value:
  triaged_at__gt = null

  # Optional string.
  # Default value:
  triaged_at__lt = null

  # Optional bool.
  # Default value:
  triaged_at__null = null

  # Optional string.
  # Default value:
  closed_at__gt = null

  # Optional string.
  # Default value:
  closed_at__lt = null

  # Optional bool.
  # Default value:
  closed_at__null = null

  # Optional string.
  # Default value:
  disclosed_at__gt = null

  # Optional string.
  # Default value:
  disclosed_at__lt = null

  # Optional bool.
  # Default value:
  disclosed_at__null = null

  # Optional bool.
  # Default value:
  reporter_agreed_on_going_public = null

  # Optional string.
  # Default value:
  bounty_awarded_at__gt = null

  # Optional string.
  # Default value:
  bounty_awarded_at__lt = null

  # Optional bool.
  # Default value:
  bounty_awarded_at__null = null

  # Optional string.
  # Default value:
  swag_awarded_at__gt = null

  # Optional string.
  # Default value:
  swag_awarded_at__lt = null

  # Optional bool.
  # Default value:
  swag_awarded_at__null = null

  # Optional string.
  # Default value:
  last_report_activity_at__gt = null

  # Optional string.
  # Default value:
  last_report_activity_at__lt = null

  # Optional string.
  # Default value:
  first_program_activity_at__gt = null

  # Optional string.
  # Default value:
  first_program_activity_at__lt = null

  # Optional bool.
  # Default value:
  first_program_activity_at__null = null

  # Optional string.
  # Default value:
  last_program_activity_at__gt = null

  # Optional string.
  # Default value:
  last_program_activity_at__lt = null

  # Optional bool.
  # Default value:
  last_program_activity_at__null = null

  # Optional string.
  # Default value:
  last_activity_at__gt = null

  # Optional string.
  # Default value:
  last_activity_at__lt = null

  # Optional string.
  # Default value:
  last_public_activity_at__gt = null

  # Optional string.
  # Default value:
  last_public_activity_at__lt = null

  # Optional string.
  # Default value:
  keyword = null

  # Optional map of string.
  # Default value:
  custom_fields = null
}
```