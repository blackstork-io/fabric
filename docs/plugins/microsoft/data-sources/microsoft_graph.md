---
title: "`microsoft_graph` data source"
plugin:
  name: blackstork/microsoft
  description: "The `microsoft_graph` data source queries Microsoft Graph"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/microsoft/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/microsoft" "microsoft" "v0.4.2" "microsoft_graph" "data source" >}}

## Description
The `microsoft_graph` data source queries Microsoft Graph.

## Installation

To use `microsoft_graph` data source, you must install the plugin `blackstork/microsoft`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/microsoft" = ">= v0.4.2"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration arguments:

```hcl
config data microsoft_graph {
  # The Azure client ID
  #
  # Required string.
  # For example:
  client_id = "some string"

  # The Azure client secret. Required if private_key_file/privat_key/cert_thumbprint is not provided.
  #
  # Optional string.
  # Default value:
  client_secret = null

  # The Azure tenant ID
  #
  # Required string.
  # For example:
  tenant_id = "some string"

  # The path to the private key file. Ignored if private_key/client_secret is provided.
  #
  # Optional string.
  # Default value:
  private_key_file = null

  # The private key contents. Ignored if client_secret is provided.
  #
  # Optional string.
  # Default value:
  private_key = null

  # The key passphrase. Ignored if client_secret is provided.
  #
  # Optional string.
  # Default value:
  key_passphrase = null
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data microsoft_graph {
  # The API version
  #
  # Optional string.
  # Default value:
  api_version = "beta"

  # The endpoint to query
  #
  # Required string.
  # For example:
  endpoint = "/security/incidents"

  # The query parameters
  #
  # Optional map of string.
  # Default value:
  query_params = null

  # Number of objects to be returned
  #
  # Optional number.
  # Must be >= 1
  # Default value:
  objects_size = 50

  # Return only the list of objects. If `false`, returns an object with `objects` and `totalCount` fields
  #
  # Optional bool.
  # Default value:
  only_objects = true

  # If API endpoint response should be treated as a list or as an object. If set to `true`, `only_objects`, `query_params` and `objects_size` are ignored.
  #
  # Optional bool.
  # Default value:
  is_object_endpoint = false
}
```