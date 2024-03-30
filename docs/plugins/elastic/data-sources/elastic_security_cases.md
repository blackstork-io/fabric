---
title: elastic_security_cases
plugin:
  name: blackstork/elastic
  description: ""
  tags: []
  version: "v0.4.0"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/elastic/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/elastic" "elastic" "v0.4.0" "elastic_security_cases" "data source" >}}

## Installation

To use `elastic_security_cases` data source, you must install the plugin `blackstork/elastic`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/elastic" = ">= v0.4.0"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data elastic_security_cases {
    api_key = <list of string>  # optional
    api_key_str = <string>  # optional
    kibana_endpoint_url = <string>  # required
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data elastic_security_cases {
    assignees = <list of string>  # optional
    default_search_operator = <string>  # optional
    from = <string>  # optional
    owner = <list of string>  # optional
    reporters = <list of string>  # optional
    search = <string>  # optional
    search_fields = <list of string>  # optional
    severity = <string>  # optional
    size = <number>  # optional
    sort_field = <string>  # optional
    sort_order = <string>  # optional
    space_id = <string>  # optional
    status = <string>  # optional
    tags = <list of string>  # optional
    to = <string>  # optional
}
```