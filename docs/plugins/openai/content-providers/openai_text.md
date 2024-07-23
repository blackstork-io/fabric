---
title: "`openai_text` content provider"
plugin:
  name: blackstork/openai
  description: ""
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/openai/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/openai" "openai" "v0.4.2" "openai_text" "content provider" >}}

## Installation

To use `openai_text` content provider, you must install the plugin `blackstork/openai`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/openai" = ">= v0.4.2"
  }
}
```

Note the version constraint set for the plugin.


#### Configuration

The content provider supports the following configuration arguments:

```hcl
config content openai_text {
  # Optional string.
  # Default value:
  system_prompt = null

  # Required string.
  # For example:
  api_key = "some string"

  # Optional string.
  # Default value:
  organization_id = null
}
```

#### Usage

The content provider supports the following execution arguments:

```hcl
content openai_text {
  # Required string.
  # For example:
  prompt = "Summarize the following text: {{.vars.text_to_summarize}}"

  # Optional string.
  # Must have a length of at least 1
  # Default value:
  model = "gpt-3.5-turbo"
}
```

