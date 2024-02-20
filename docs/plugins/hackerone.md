---
title: blackstork/hackerone
weight: 20
type: docs
---

# `blackstork/hackerone` plugin

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/hackerone" = "=> v0.0.0-dev"
  }
}
```

## Data sources

The plugin has the following data sources available:

### `hackerone_reports`

#### Configuration

The data source supports the following configuration parameters:

```hcl
config data hackerone_reports {
    api_token = <string>  # required
    api_username = <string>  # required
}
```

#### Usage

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