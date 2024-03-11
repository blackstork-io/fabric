---
title: virustotal_api_usage
plugin:
  name: blackstork/virustotal
  description: ""
  tags: []
  version: "v0.4.0"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/virustotal/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/virustotal" "virustotal" "v0.4.0" "virustotal_api_usage" "data source" >}}

## Installation

To use `virustotal_api_usage` data source, you must install the plugin `blackstork/virustotal`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/virustotal" = ">= v0.4.0"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data virustotal_api_usage {
    api_key = <string>  # required
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data virustotal_api_usage {
    end_date = <string>  # optional
    group_id = <string>  # optional
    start_date = <string>  # optional
    user_id = <string>  # optional
}
```