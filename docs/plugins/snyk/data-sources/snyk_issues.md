---
title: snyk_issues
plugin:
  name: blackstork/snyk
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/snyk/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/snyk" "snyk" "v0.4.1" "snyk_issues" "data source" >}}

## Installation

To use `snyk_issues` data source, you must install the plugin `blackstork/snyk`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/snyk" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data snyk_issues {
    api_key = <string>  # required
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data snyk_issues {
    created_after = <string>  # optional
    created_before = <string>  # optional
    effective_severity_level = <list of string>  # optional
    group_id = <string>  # optional
    ignored = <bool>  # optional
    limit = <number>  # optional
    org_id = <string>  # optional
    scan_item_id = <string>  # optional
    scan_item_type = <string>  # optional
    status = <list of string>  # optional
    type = <string>  # optional
    updated_after = <string>  # optional
    updated_before = <string>  # optional
}
```