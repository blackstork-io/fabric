---
title: "`microsoft_security` data source"
plugin:
  name: blackstork/microsoft
  description: "The `microsoft_security` data source queries Microsoft Security API"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/microsoft/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/microsoft" "microsoft" "v0.4.2" "microsoft_security" "data source" >}}

## Description
The `microsoft_security` data source queries Microsoft Security API.

## Installation

To use `microsoft_security` data source, you must install the plugin `blackstork/microsoft`.

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
config data microsoft_security {
  # The Azure client ID
  #
  # Required string.
  # For example:
  client_id = "some string"

  # The Azure client secret. Required if `private_key_file` or `private_key` is not provided.
  #
  # Optional string.
  # Default value:
  client_secret = null

  # The Azure tenant ID
  #
  # Required string.
  # For example:
  tenant_id = "some string"

  # The path to the private key file. Ignored if `private_key` or `client_secret` is provided.
  #
  # Optional string.
  # Default value:
  private_key_file = null

  # The private key contents. Ignored if `client_secret` is provided.
  #
  # Optional string.
  # Default value:
  private_key = null

  # The key passphrase. Ignored if `client_secret` is provided.
  #
  # Optional string.
  # Default value:
  key_passphrase = null
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data microsoft_security {
  # API endpoint to query
  #
  # Required string.
  # For example:
  endpoint = "/users"

  # HTTP query parameters
  #
  # Optional map of string.
  # Default value:
  query_params = null

  # Number of objects to be returned
  #
  # Optional number.
  # Must be >= 1
  # Default value:
  size = 50

  # Indicates if API endpoint serves a single object. If set to `true`, `query_params` and `size` arguments are ignored.
  #
  # Optional bool.
  # Default value:
  is_object_endpoint = false
}
```