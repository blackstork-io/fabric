---
title: "`microsoft_security_query` data source"
plugin:
  name: blackstork/microsoft
  description: "The `microsoft_defender_query` data source queries Microsoft Security API"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/microsoft/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/microsoft" "microsoft" "v0.4.2" "microsoft_security_query" "data source" >}}

## Description
The `microsoft_defender_query` data source queries Microsoft Security API.

## Installation

To use `microsoft_security_query` data source, you must install the plugin `blackstork/microsoft`.

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
config data microsoft_security_query {
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
data microsoft_security_query {
  # Advanced hunting query to run
  #
  # Required string.
  # For example:
  query = "DeviceRegistryEvents | where Timestamp >= ago(30d) | where isnotempty(RegistryKey) and isnotempty(RegistryValueName) | limit 5"
}
```