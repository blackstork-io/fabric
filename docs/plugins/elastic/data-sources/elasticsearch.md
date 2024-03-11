---
title: elasticsearch
plugin:
  name: blackstork/elastic
  description: ""
  tags: []
  version: "v0.0.0-dev"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/elastic/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/elastic" "elastic" "v0.0.0-dev" "elasticsearch" "data source" >}}

## Installation

To use `elasticsearch` data source, you must install the plugin `blackstork/elastic`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/elastic" = ">= v0.0.0-dev"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data elasticsearch {
    api_key = <list of string>  # optional
    api_key_str = <string>  # optional
    base_url = <string>  # optional
    basic_auth_password = <string>  # optional
    basic_auth_username = <string>  # optional
    bearer_auth = <string>  # optional
    ca_certs = <string>  # optional
    cloud_id = <string>  # optional
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data elasticsearch {
    aggs = <map of dynamic>  # optional
    fields = <list of string>  # optional
    id = <string>  # optional
    index = <string>  # required
    only_hits = <bool>  # optional
    query = <map of dynamic>  # optional
    query_string = <string>  # optional
    size = <number>  # optional
}
```