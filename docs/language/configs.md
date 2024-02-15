---
title: Configuration
type: docs
weight: 20
---

# Configuration

## Global configuration

The `fabric` configuration block defines a global configuration for Fabric. It's used for defining plugin dependencies, the paths to local directories, etc.

```hcl
fabric {
  ...
}
```

There can be only one `fabric` block defined within the codebase.

### Supported arguments

- `plugin_versions`: (required) a map that matches name-spaced plugin names to the version constraints in SemVer (see Terraform [version constraint syntax](https://developer.hashicorp.com/terraform/language/expressions/version-constraints#version-constraint-syntax)).
- `cache_dir`: (optional) a path to a directory on the local file system. The default value is `.fabric` directory in the current folder. If the directory doesn't exist, Fabric will create it on the first run.

### Supported nested blocks

- `plugin_registry` – (optional) a block that defines available plugin registries. At the moment, the block accepts only one attribute:

  ```hcl
  plugin_registry {
    mirror_dir = "<path>"
  }
  ```

  - `mirror_dir` – (optional) a path to a directory on the local filesystem with plugin binaries.

### Example

```hcl
fabric {

  cache_dir = "./.fabric"

  plugins_registry {
    mirror_dir = "/tmp/local-mirror/plugins"
  }

  plugin_versions = {
    "blackstork/elasticsearch" = "1.2.3"
    "blackstork/openai" = "=11.22.33"
  }
}
```

## Data source configuration

A Fabric plugin can include one or more data sources. For example, `blackstork/github` plugin includes `github_issues` data source.

A data source loads data from a local or an external data store, platform, and service.

Data sources are configured using `config` block:

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
