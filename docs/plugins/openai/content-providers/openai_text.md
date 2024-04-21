---
title: openai_text
plugin:
  name: blackstork/openai
  description: ""
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/openai/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/openai" "openai" "v0.4.1" "openai_text" "content provider" >}}

## Installation

To use `openai_text` content provider, you must install the plugin `blackstork/openai`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/openai" = ">= v0.4.1"
  }
}
```

Note the version constraint set for the plugin.


#### Configuration

The content provider supports the following configuration parameters:

```hcl
config content openai_text {
  # Optional string. Default value:
  system_prompt = null

  # Required string. For example:
  api_key = "some string"

  # Optional string. Default value:
  organization_id = null
}
```

#### Usage

The content provider supports the following execution parameters:

```hcl
content openai_text {
  # Required string. For example:
  prompt = "some string"

  # Optional string. Default value:
  model = null
}
```

