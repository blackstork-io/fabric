---
title: hackerone_reports
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

The data source supports the following configuration parameters:

```hcl
config "data" "hackerone_reports" {
  # Required. For example:
  api_username = "some string"

  # Required. For example:
  api_token = "some string"
}

```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data "hackerone_reports" {
  # Optional. Default value:
  size = null

  # Optional. Default value:
  page_number = null

  # Optional. Default value:
  sort = null

  # Optional. Default value:
  program = null

  # Optional. Default value:
  inbox_ids = null

  # Optional. Default value:
  reporter = null

  # Optional. Default value:
  assignee = null

  # Optional. Default value:
  state = null

  # Optional. Default value:
  id = null

  # Optional. Default value:
  weakness_id = null

  # Optional. Default value:
  severity = null

  # Optional. Default value:
  hacker_published = null

  # Optional. Default value:
  created_at__gt = null

  # Optional. Default value:
  created_at__lt = null

  # Optional. Default value:
  submitted_at__gt = null

  # Optional. Default value:
  submitted_at__lt = null

  # Optional. Default value:
  triaged_at__gt = null

  # Optional. Default value:
  triaged_at__lt = null

  # Optional. Default value:
  triaged_at__null = null

  # Optional. Default value:
  closed_at__gt = null

  # Optional. Default value:
  closed_at__lt = null

  # Optional. Default value:
  closed_at__null = null

  # Optional. Default value:
  disclosed_at__gt = null

  # Optional. Default value:
  disclosed_at__lt = null

  # Optional. Default value:
  disclosed_at__null = null

  # Optional. Default value:
  reporter_agreed_on_going_public = null

  # Optional. Default value:
  bounty_awarded_at__gt = null

  # Optional. Default value:
  bounty_awarded_at__lt = null

  # Optional. Default value:
  bounty_awarded_at__null = null

  # Optional. Default value:
  swag_awarded_at__gt = null

  # Optional. Default value:
  swag_awarded_at__lt = null

  # Optional. Default value:
  swag_awarded_at__null = null

  # Optional. Default value:
  last_report_activity_at__gt = null

  # Optional. Default value:
  last_report_activity_at__lt = null

  # Optional. Default value:
  first_program_activity_at__gt = null

  # Optional. Default value:
  first_program_activity_at__lt = null

  # Optional. Default value:
  first_program_activity_at__null = null

  # Optional. Default value:
  last_program_activity_at__gt = null

  # Optional. Default value:
  last_program_activity_at__lt = null

  # Optional. Default value:
  last_program_activity_at__null = null

  # Optional. Default value:
  last_activity_at__gt = null

  # Optional. Default value:
  last_activity_at__lt = null

  # Optional. Default value:
  last_public_activity_at__gt = null

  # Optional. Default value:
  last_public_activity_at__lt = null

  # Optional. Default value:
  keyword = null

  # Optional. Default value:
  custom_fields = null
}

```