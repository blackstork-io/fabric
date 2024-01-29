---
title: Data blocks
type: docs
weight: 1
---

# Data blocks

`data` block defines a call to a data plugin. The data plugins provide data from external sources for content rendering. The order of the `data` block definitions in the document does not matter.

```hcl
data <plugin-name> "<result-name>" {
  ...
}

document "foobar" {

  data <plugin-name> "<result-name>" {
    ...
  }

}
```

If `data` block is defined on a root level of the configuration file, both names - the name of the plugin and the result name - must be provided. The plugin name / result name pair must be unique within the codebase, as it will be used as an identifier when referencing this block. The data blocks defined outside the document are not executed independently but must be referenced inside the document template.

If the `data` block is defined in the document, it must be on the root level of the document. The name of the plugin and the name of the result are required. The plugin name / result name pair must be unique in the scope of the document, since it will be used as an identifier.

The arguments set in the block are the input parameters for the data plugin, together with the plugin configuration. The data returned by the plugin is set in the global context map under a path `data.<plugin-name>.<result-name>`


## Supported Arguments

- `config` – _(optional)_ a reference to a named config block defined on a root level. If provided, it takes precedence over the default config for the plugin.

The plugin might define other supported arguments.


## Supported Nested Blocks

- `meta`
- `config` – _(optional)_ an inline config block for the plugin. If provided, it takes precedence over the default configuration for the plugin.

Other nested blocks are not supported.


## References

If the label `ref` is used instead of `<plugin-name>`, the block references another data block defined on a root level. The name of the referer block is optional if the block is defined within the document. If the referer block (with label `ref`) is defined on a root level of the config file, the name is required.

{{< hint warning >}}
#### Overriding existing blocks
The data ref block with the already used name (explicit or inherited) will override the previously defined block.  
{{< /hint >}}

```hcl
data elasticsearch "foo" {
  ...
}

document "overview" {

  data ref {
    base = data.elasticsearch.foo
    ...
  }

  data ref "foo2" {
    base = data.elasticsearch.foo
    ...
  }

}
```

Every referer block must have a `base` attribute set, pointing to a named data block defined on a root level in the config file.

Other supported arguments are defined by the plugin of the referent block (`elasticsearch` in the above example). If any arguments are provided in the referer block, they take precedence over the arguments set in the referent block.

Referer blocks can not contain nested blocks.
