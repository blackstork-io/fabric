---
title: blackstork/openai
weight: 20
type: docs
---

# `blackstork/openai` plugin

## Installation

To install the plugin, add it to `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), with a version constraint restricting which available versions of the plugin the codebase is compatible with:

```hcl
fabric {
  plugin_versions = {
    "blackstork/openai" = "=> v0.0.0-dev"
  }
}
```



## Content providers
The plugin has the following content providers available:

### `openai_text`

#### Configuration

The content provider supports the following configuration parameters:

```hcl
config content openai_text {
    api_key = <string>  # required
    organization_id = <string>  # optional
    system_prompt = <string>  # optional
}
```

#### Usage

The content source supports the following parameters in the content blocks:

```hcl
content openai_text {
    model = <string>  # optional
    prompt = <string>  # required
}
```
