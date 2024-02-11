---
title: Configuration
type: docs
weight: 20
---

# Configuration

## Global Configuration

Configuration block type `fabric` is used to configure Fabric. It can define locations of the local directories and plugin dependencies.

```hcl
fabric {
  ...
}
```

There can be only one `fabric` block defined within the codebase.

### Supported Arguments

- `plugin_versions`: (required) a map that matches namespaced plugin names to the version constraints in SemVer (see Terraform's [version constraint syntax](https://developer.hashicorp.com/terraform/language/expressions/version-constraints#version-constraint-syntax)).
- `cache_dir`: (optional) a path to a directory on the local file system. The default value is `./.fabric`. If the directory does not exist, it will be created on the first run of `fabric`.

### Supported Nested Blocks

- `plugin_registry` – (optional) a block that defines available plugin registries. The block accepts only one attribute

  ```hcl
  plugin_registry {
    cache_dir = "/tmp/plugins/"
  }
  ```

  - `mirror_dir` – (optional) a path to a directory on the local filesystem with plugin archives.

### Example

```hcl
fabric {

  plugins_registry {
    mirror_dir = "/tmp/plugins/"
  }

  cache_dir = "./.fabric"

  plugin_versions = {
    "blackstork/data.elasticsearch" = "1.2.3"
    "blackstork/content.openai" = "=11.22.33"
  }
}
```

## Plugin Configuration

`config` block defines a configuration for a plugin:

```hcl

config <plugin-type> <plugin-name> {
  ...
}

config <plugin-type> <plugin-name> "<config-name>" {
  ...
}
```

A plugin type matches a block type it powers: `<plugin-type>` is either `content` or `data`. Both plugin type and plugin name are required.

If `<config-name>` is not provided, the block is treated as a default configuration for the plugin of a specified type and name. Every time the plugin is executed (while running `content` or `data` blocks), the configuration will be passed to the plugin as one of the arguments.

If `<config-name>` is set, the config block can be explicitely referenced inside the `content` or `data` block. This is helpful if there is a need to have multiple configurations for the same plugin.

### Supported Arguments

The arguments that are allowed in the configuration block are plugin-specific – every plugin defines the configuration options supported. See [Plugins]({{< ref "plugins.md" >}}) for the details on the plugin configuation paramters.

### Supported Nested Blocks

Nested blocks are not supported inside the `config` blocks.

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
