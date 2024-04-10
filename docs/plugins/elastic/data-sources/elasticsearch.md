---
title: elasticsearch
plugin:
  name: blackstork/elastic
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/elastic/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/elastic" "elastic" "v0.4.1" "elasticsearch" "data source" >}}

## Installation

To use `elasticsearch` data source, you must install the plugin `blackstork/elastic`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/elastic" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config "data" "elasticsearch" {
  # Optional. Default value:
  base_url = null

  # Optional. Default value:
  cloud_id = null

  # Optional. Default value:
  api_key_str = null

  # Optional. Default value:
  api_key = null

  # Optional. Default value:
  basic_auth_username = null

  # Optional. Default value:
  basic_auth_password = null

  # Optional. Default value:
  bearer_auth = null

  # Optional. Default value:
  ca_certs = null
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data "elasticsearch" {
  # Required. For example:
  index = "some string"

  # Optional. Default value:
  id = null

  # Optional. Default value:
  query_string = null

  # Optional. Default value:
  query = null

  # Optional. Default value:
  aggs = null

  # Optional. Default value:
  only_hits = null

  # Optional. Default value:
  fields = null

  # Optional. Default value:
  size = null
}
```