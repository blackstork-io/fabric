---
title: Configuration
type: docs
weight: 2
---

# Global configuration ([#5](https://github.com/blackstork-io/fabric/issues/5))

Fabric can be configured through a global configuration block:
```hcl
fabric {
  ...
}
```

There can be only one `fabric` block defined within the codebase.

## Supported Arguments

- `cache_dir` – _(optional)_ a path to a directory on the local FS. The default value is `./.fabric`. If the directory does not exist, it will be created on the first run of `fabric`.
- `plugin_versions` – (required) a map that matches namespaced plugin names to the version constraints (SemVer, in Terraform's [version constraint syntax](https://developer.hashicorp.com/terraform/language/expressions/version-constraints#version-constraint-syntax))

No other arguments are supported.

## Supported Nested Blocks

- `plugin_registry` – _(optional)_ block that defines available plugin registries. The block accepts only one attribute
  ```hcl
  plugin_registry {
    mirror_dir = "/tmp/plugins/"
  }
  ```
  - `mirror_dir` – _(optional)_ a path to a directory on the local FS with plugin archives.


# Plugin configurations ([#4](https://github.com/blackstork-io/fabric/issues/4))

`config` block defines a configuration for a plugin.

```hcl

config <plugin-type> <plugin-name> {
    ...
}

config <plugin-type> <plugin-name> "<config-name>" {
    ...
}
```

`<plugin-type>` is either `content` or `data`.

If `<config-name>` is not provided, the block is treated as a default configuration for the plugin of a specified type (`content` or `data`) with a specified name. Every time the plugin is executed (during the execution of `content` or `data` blocks), the configuration will be passed to the plugin.

If `<config-name>` is set, the config block can be explicitely referenced inside the `content` or `data` block. This is helpful if there is a need to have multiple configurations for the same plugin.


## Supported Arguments
The arguments supported in the block are plugin-specific – every plugin defines the configuration options supported.


## Supported Nested Blocks

Nested blocks are not supported inside `config` blocks.

