---
title: Configuration
description: Learn how to configure Fabric, data sources and content providers.
type: docs
weight: 20
---

# Configuration

## Global configuration

The `fabric` configuration block serves as the global configuration for Fabric, offering a centralized space for defining essential aspects such as plugin dependencies and local directory paths.

```hcl
fabric {
  ...
}
```

Within the codebase, only one `fabric` block can be defined.

### Supported arguments

- `plugin_versions`: (required) a map that aligns namespaced plugin names with version constraints in SemVer (refer to Terraform [version constraint syntax](https://developer.hashicorp.com/terraform/language/expressions/version-constraints#version-constraint-syntax)).
- `cache_dir`: (optional) a path to a directory on the local file system. The default value is `.fabric` directory in the current folder. If the directory doesn't exist, Fabric creates it upon the first run.

### Supported nested blocks

- `plugin_registry`: (optional) a block defines available plugin registries and can include the following arguments:

  ```hcl
  plugin_registry {
    base_url = "<url>"
    mirror_dir = "<path>"
  }
  ```

  - `base_url`: (optional) the base URL of the plugin registry. Default value: `https://registry.blackstork.io`
  - `mirror_dir`: (optional) the path to a directory on the local filesystem containing plugin binaries.

### Example

```hcl
fabric {

  cache_dir = "./.fabric"

  plugins_registry {
    mirror_dir = "/tmp/local-mirror/plugins"
  }

  plugin_versions = {
    "blackstork/elastic" = "1.2.3"
    "blackstork/openai" = "=11.22.33"
  }
}
```

## Data source configuration

A Fabric plugin can include one or more data sources. For example, `blackstork/github` plugin includes `github_issues` data source.

A data source within Fabric serves the purpose of loading data from either a local or an external data store, platform, or service.

Configuration for data sources is accomplished through the use of the `config` block:

```hcl
config data <data-source-name> "<name>" {
  ...
}
```

If the block is named (`<name>` is provided), the `config` block can be referenced in a `config` argument inside `data` blocks. This is helpful if there is a need to have more than one configuration for the same data source.

If `<name>` isn't provided, the configuration acts as a default configuration for a specified data source.

### Supported arguments

The arguments allowed in the configuration block depend on the data source. See [Plugins]({{< ref "plugins.md" >}}) for the details on the configuration parameters supported.

### Supported nested blocks

Nested blocks aren't supported inside the `config` blocks.

### Example

```hcl
config data csv {
  delimiter = ";"
}

data csv "events_a" {
  path = "/tmp/events-a.csv"
}

document "test-document" {

  data ref {
    base = data.csv.events_a
  }

  data csv "events_b" {
    config {
      delimiter = ",";
    }

    path = "/tmp/events-b.csv"
  }
}
```

## Content provider configuration

A Fabric plugin can include one or more content providers. For example, `blackstork/openai` plugin includes `openai_text` content provider.
Content providers generate Markdown content that Fabric will include into the rendered document. The provider can either render content locally or use an external API (like OpenAI API).

Content providers can be configured using `config` block:

```hcl
config content <content-provider-name> "<name>" {
  ...
}
```

If the block is named (`<name>` is provided), the `config` block can be explicitly referenced in a `config` argument inside `content`. This is helpful if there is a need to have more than one configuration for the same content provider.

If `<name>` isn't provided, the configuration acts as a default configuration for a specified content provider.

### Supported arguments

The arguments allowed in the configuration block depend on the content provider. See [Plugins]({{< ref "plugins.md" >}}) for the details on the configuration parameters supported.

### Supported nested blocks

Nested blocks aren't supported inside the `config` blocks.

### Example

```hcl
config content openai_text {
  api_key = 'some-openai-api-key'

  system_prompt = 'You are the best at saying Hi!'
}

document "test-document" {

  content openai_text {
    prompt = "Say hi!"
  }
}
```

## Next steps

See [Documents]({{< ref "documents.md" >}}) to learn how to build document templates in Fabric configuration language.
