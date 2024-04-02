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
config data hackerone_reports {
    api_token = <string>  # required
    api_username = <string>  # required
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data hackerone_reports {
    assignee = <list of string>  # optional
    bounty_awarded_at__gt = <string>  # optional
    bounty_awarded_at__lt = <string>  # optional
    bounty_awarded_at__null = <bool>  # optional
    closed_at__gt = <string>  # optional
    closed_at__lt = <string>  # optional
    closed_at__null = <bool>  # optional
    created_at__gt = <string>  # optional
    created_at__lt = <string>  # optional
    custom_fields = <map of string>  # optional
    disclosed_at__gt = <string>  # optional
    disclosed_at__lt = <string>  # optional
    disclosed_at__null = <bool>  # optional
    first_program_activity_at__gt = <string>  # optional
    first_program_activity_at__lt = <string>  # optional
    first_program_activity_at__null = <bool>  # optional
    hacker_published = <bool>  # optional
    id = <list of number>  # optional
    inbox_ids = <list of number>  # optional
    keyword = <string>  # optional
    last_activity_at__gt = <string>  # optional
    last_activity_at__lt = <string>  # optional
    last_program_activity_at__gt = <string>  # optional
    last_program_activity_at__lt = <string>  # optional
    last_program_activity_at__null = <bool>  # optional
    last_public_activity_at__gt = <string>  # optional
    last_public_activity_at__lt = <string>  # optional
    last_report_activity_at__gt = <string>  # optional
    last_report_activity_at__lt = <string>  # optional
    page_number = <number>  # optional
    program = <list of string>  # optional
    reporter = <list of string>  # optional
    reporter_agreed_on_going_public = <bool>  # optional
    severity = <list of string>  # optional
    size = <number>  # optional
    sort = <string>  # optional
    state = <list of string>  # optional
    submitted_at__gt = <string>  # optional
    submitted_at__lt = <string>  # optional
    swag_awarded_at__gt = <string>  # optional
    swag_awarded_at__lt = <string>  # optional
    swag_awarded_at__null = <bool>  # optional
    triaged_at__gt = <string>  # optional
    triaged_at__lt = <string>  # optional
    triaged_at__null = <bool>  # optional
    weakness_id = <list of number>  # optional
}
```