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
    api_key = <string>  # optional
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data nist_nvd_cves {
    cpe_name = <string>  # optional
    cve_id = <string>  # optional
    cvss_v3_metrics = <string>  # optional
    cvss_v3_severity = <string>  # optional
    cwe_id = <string>  # optional
    has_cert_alerts = <bool>  # optional
    has_cert_notes = <bool>  # optional
    has_kev = <bool>  # optional
    is_vulnerable = <bool>  # optional
    keyword_exact_match = <bool>  # optional
    keyword_search = <string>  # optional
    last_mod_end_date = <string>  # optional
    last_mod_start_date = <string>  # optional
    limit = <number>  # optional
    no_rejected = <bool>  # optional
    pub_end_date = <string>  # optional
    pub_start_date = <string>  # optional
    source_identifier = <string>  # optional
    virtual_match_string = <string>  # optional
}
```