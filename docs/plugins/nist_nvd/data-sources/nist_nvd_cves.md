---
title: nist_nvd_cves
plugin:
  name: blackstork/nist_nvd
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/nistnvd/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/nist_nvd" "nist_nvd" "v0.4.1" "nist_nvd_cves" "data source" >}}

## Installation

To use `nist_nvd_cves` data source, you must install the plugin `blackstork/nist_nvd`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/nist_nvd" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data nist_nvd_cves {
  # Optional string.
  # Default value:
  api_key = null
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data nist_nvd_cves {
  # Optional string.
  # Default value:
  last_mod_start_date = null

  # Optional string.
  # Default value:
  last_mod_end_date = null

  # Optional string.
  # Default value:
  pub_start_date = null

  # Optional string.
  # Default value:
  pub_end_date = null

  # Optional string.
  # Default value:
  cpe_name = null

  # Optional string.
  # Default value:
  cve_id = null

  # Optional string.
  # Default value:
  cvss_v3_metrics = null

  # Optional string.
  # Default value:
  cvss_v3_severity = null

  # Optional string.
  # Default value:
  cwe_id = null

  # Optional string.
  # Default value:
  keyword_search = null

  # Optional string.
  # Default value:
  virtual_match_string = null

  # Optional string.
  # Default value:
  source_identifier = null

  # Optional bool.
  # Default value:
  has_cert_alerts = null

  # Optional bool.
  # Default value:
  has_kev = null

  # Optional bool.
  # Default value:
  has_cert_notes = null

  # Optional bool.
  # Default value:
  is_vulnerable = null

  # Optional bool.
  # Default value:
  keyword_exact_match = null

  # Optional bool.
  # Default value:
  no_rejected = null

  # Optional number.
  # Default value:
  limit = null
}
```