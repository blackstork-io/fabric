---
title: Configuration
description: Learn how to configure Fabric, data sources and content providers.
type: docs
weight: 20
---

# Configuration

## Global configuration

The `fabric` configuration block serves as the global configuration for Fabric. Use it to define
configuration properties, including default behavior, plugin dependencies, and local paths.

```hcl
fabric {
  plugin_versions = {
    "blackstork/elastic" = ">0.4.1"
    "blackstork/openai" = ">0.4.1"
  }
}
```

Within the codebase, only one `fabric` block can be defined.

### Supported arguments

- `plugin_versions`: (optional) a map of plugin dependencies. The version constraints are defined in
  SemVer (refer to Terraform [version constraint
  syntax](https://developer.hashicorp.com/terraform/language/expressions/version-constraints#version-constraint-syntax)).
- `expose_env_vars_with_pattern`: (optional) a glob pattern for environment variable names. Only the
  environment variables that match the pattern are exposed in the evaluation context. Default value
  is `FABRIC_*`. To expose all environment variables, set the value to `*`. See [Environment
  variables]({{< ref "configs.md#environment-variables">}}) for the details.
- `cache_dir`: (optional) a path to a directory on the local file system. The default value is
  `.fabric` - a directory in the current folder. If the directory doesn't exist, Fabric will create
  it during the first run.

To install all dependencies defined in `plugin_versions`, run `fabric install` command (see
[Installing plugins]({{< ref "install.md#installing-plugins" >}}) for more details)

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

## Block configuration

The data sources, content provides and publishers can be configured using `config` blocks.

`config` block arguments are specific to a data source, content provider or a publisher the block if
configuring.

`config` block must be defined on a root level of Fabric file, outside `document` block. The
signature of `config` block consists of a block type selector (`content`, `data` or `publish`), a
data source / content provider / publisher name, and a block name:

```hcl
config <block-type> <source/provider/publisher-name> "<name>" {
  # ...
}
```

If `<name>` isn't provided, the configuration acts as a default configuration for a specified data
source / content provider / publisher.

If the block has a name (`<name>` is specified), `config` block can be referenced in a `config` argument.
This is helpful if there is a need to have more than one configuration available.

### Supported arguments

The arguments allowed in the configuration block depend on the data source / content provider / publisher. See the documentation for [Data Sources]({{< ref data-sources.md >}}), [Content Providers]({{< ref content-providers.md >}}), and [Publishers]({{< ref publishers.md >}}) for the details on the configuration parameters supported.

### Supported nested blocks

- `meta`: (optional) a block containing metadata for the block. See [Metadata]({{< ref "configs.md#metadata" >}}) for details.

### Example

```hcl
config data csv {
  delimiter = ";"
}

config content openai_text {
  api_key = env.OPENAI_API_KEY
  system_prompt = "You are the best at saying Hi!"
}

document "test-document" {

  data csv "events_a" {
    path = "/tmp/events-a.csv"
  }

  data csv "events_b" {
    # Overriding the default configuration for CSV data source
    config {
      delimiter = ","
    }

    path = "/tmp/events-b.csv"
  }

  content openai_text {
    prompt = "Say hi!"
  }
}
```

## Environment variables

Fabric templates can be configured with environment variables, either set in the shell or provided
in `.env` file.

The values of environment variables are available through global `env` object and in the context
under `.env` root.

For example:

```hcl
config data elasticsearch {
  basic_auth_username = "elastic"
  basic_auth_password = env.ELASTICSEARCH_PASSWORD
}

content text {
  local_var = query_jq(".env.FOOBAR | split(\",\") | length")
  value = "There are {{ .vars.local }} elements in `FOOBAR` env var"
}
```

{{< hint warning >}}

In accordance with the principle of least privilege, Fabric limits which environment variables are exposed to plugins in the evaluation context.

Only the variables with names that match the pattern set in `expose_env_vars_with_pattern` argument
(see [Global configuration]({{< ref "configs.md#global-configuration">}})) are available under `.env` key in the evaluation context.

It's recommended to use the most restrictive glob pattern possible.

Note, that the filtering doesn't apply to the variables exposed in `env` object: it's the user's responsibility to make sure that their Fabric template doesn't expose sensitive environment variables.
{{< /hint >}}

## Metadata

Define the metadata for a template or a block using the `meta` block. The metadata includes the
template name, tags, authors' names, license, version, and other relevant information.

`meta` block supports following arguments:

- `name` — name of the template or the block
- `description` — description of the template or the block
- `url` — a string with a URL
- `license` — the license that applies to the template or the block
- `authors` — a list of the authors
- `tags` — a list of tags
- `updated_at` — ISO8601-formatted date with time
- `version` — version of the template or the block

For example:

```hcl
# Document template with metadata
document "mitre_ctid_campaign_report" {

  meta {
    name = "MITRE CTID Campaign Report Template"

    description = <<-EOT
      The Campaign Report is designed to highlight new information related to a threat actor or
      capabilities. This should focus on new information and highlight how it poses a changed risk
      to your organization. This should not be an exhaustive product cataloguing all information
      about the topic, but rather a succinct report designed to convey a change in the status quo to
      the intended recipient.
    EOT

    url = "https://github.com/center-for-threat-informed-defense/cti-blueprints"

    license = "Apache License 2.0"
    tags = ["mitre", "campaign"]

    updated_at = "2024-01-22T10:00:01+01:00"
  }

}

content text "disclaimer" {
  meta {
    name = "Disclaimer text"
    tags = ["foo", "bar"]
  }

  value = "Some disclaimer text"
}
```

## Next steps

See [Documents]({{< ref "documents.md" >}}) to learn how to build document templates in Fabric configuration language.
