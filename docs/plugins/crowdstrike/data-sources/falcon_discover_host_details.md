---
title: "`falcon_discover_host_details` data source"
plugin:
  name: blackstork/crowdstrike
  description: "The `falcon_discover_host_details` data source fetches host details from Falcon Discover Host API"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/crowdstrike/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/crowdstrike" "crowdstrike" "v0.4.2" "falcon_discover_host_details" "data source" >}}

## Description
The `falcon_discover_host_details` data source fetches host details from Falcon Discover Host API.

## Installation

To use `falcon_discover_host_details` data source, you must install the plugin `blackstork/crowdstrike`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/crowdstrike" = ">= v0.4.2"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration arguments:

```hcl
config data falcon_discover_host_details {
  # Client ID for accessing CrowdStrike Falcon Platform
  #
  # Required string.
  # Must be non-empty
  # For example:
  client_id = "some string"

  # Client Secret for accessing CrowdStrike Falcon Platform
  #
  # Required string.
  # Must be non-empty
  # For example:
  client_secret = "some string"

  # Member CID for MSSP
  #
  # Optional string.
  # Default value:
  member_cid = null

  # Falcon cloud abbreviation
  #
  # Optional string.
  # Must be one of: "autodiscover", "us-1", "us-2", "eu-1", "us-gov-1", "gov1"
  # For example:
  # client_cloud = "us-1"
  # 
  # Default value:
  client_cloud = null
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data falcon_discover_host_details {
  # limit the number of queried items
  #
  # Required integer.
  # For example:
  size = 42

  # Host search expression using Falcon Query Language (FQL)
  #
  # Optional string.
  # Default value:
  filter = null
}
```