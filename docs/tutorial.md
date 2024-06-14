---
title: Tutorial
description: Dive into Fabric tutorial to learn everything you need to know about using Fabric effectively. From the basics of FCL commands to advanced features such as data configurations, our tutorial provides clear, step-by-step instructions for building Fabric templates. Start improving your workflow with Fabric today.
type: docs
weight: 30
code_blocks_no_wrap: true
---

# Tutorial

This tutorial provides comprehensive guidance on using Fabric and the [Fabric Configuration Language]({{< ref "language" >}}) (FCL) for document generation. We'll systematically cover creating a basic template, incorporating data blocks, applying data filtering and mutation, installing plugins, and rendering text with external content providers.

## Prerequisites

To effectively follow this tutorial, ensure you have the following:

- Fabric CLI [installed]({{< ref "install.md" >}}) and `fabric` CLI command available
- (optional) OpenAI API token

Throughout this tutorial, the command examples were executed in macOS Sonoma `zsh` shell.

## Hello, Fabric

Let's start with a straightforward "Hello, Fabric!" template to confirm that everything is
configured correctly.

Create a new `hello.fabric` file and define a simple template:

```hcl
document "greeting" {

  content text {
    value = "Hello, Fabric!"
  }

}
```

In this code snippet, `document.greeting` block defines a template with a single anonymous content
block containing the static text "Hello, Fabric!"

To render the document, execute `fabric` command in the directory with `hello.fabric` file, or
explicitly specify a path to another directory with `--source-dir` CLI argument:

```shell
fabric render document.greeting
```

The output should resemble the following:

```shell
$ fabric render document.greeting
Hello, Fabric!
```

## Document title

Documents typically include titles, so the document block supports `title` argument as an
easy way to set a title for a document.

With the new `title` argument, `document.greeting` template would look like this:

```hcl
document "greeting" {

  title = "The Greeting"

  content text {
    value = "Hello, Fabric!"
  }

}
```

{{< hint note >}}
`title` argument for `document` block is a syntactic sugar that Fabric will translate into `content.title` block:

```hcl
content title {
  value = "The Greeting"
}
```

See [`content.title`]({{< ref "plugins/builtin/content-providers/title" >}}) content providers documentations for the details.
{{< /hint >}}

The rendered output should now include the document title:

```markdown
$ fabric render document.greeting
# The Greeting

Hello, Fabric!
```

## Variables

For this tutorial, we will use variables instead of defining the data requirements.

Change the template in the `hello.fabric` file to include a `data` block and add another
`content.text` block:

```hcl
document "greeting" {

  vars {
    solar_system = {
      planets = [
        "Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune"
      ]
      moons_count = 146
    }
  }

  title = "The Greeting"

  content text {
    value = "Hello, Fabric!"
  }

  content text {
    local_var = query_jq(".vars.solar_system.planets | length")

    value = <<-EOT
      There are {{ .vars.local }} planets and {{ .vars.solar_system.moons_count }} moons in our solar system.
    EOT
  }

}
```

The content blocks can access and transform the data available in the context (see [Evaluation Context]({{< ref context.md >}})) with [JQ queries](https://jqlang.github.io/jq/manual/). By
specifying `local_var` argument, we're using a shortcut that defines a variable called `local` for
us. The query will be executed against the context and the results will be stored in the context
under `.vars.local` path.

As you can see, `value` argument in the new content block contains a template string –
`content.text` blocks support [Go templates](https://pkg.go.dev/text/template) out-of-the-box. The
templates can access the evaluation context, so it's easy to include `local` or
`solar_system.moons_count` variable values.

The rendered output will now include the new sentence:

```shell
$ fabric render document.greeting
# The Greeting

Hello, Fabric!

There are 8 planets and 146 moons in our solar system.
```

## Content providers

Fabric seamlessly integrates with external APIs for content generation. An excellent example is the
use of the OpenAI API to dynamically generate text through prompts.

In scenarios where providing the exact text or a template string for the content block proves
challenging or impossible, leveraging generative AI for summarization becomes invaluable. This
enables users to dynamically create context-aware text.

In this tutorial, we will use the [`openai_text`]({{< ref "plugins/openai/content-providers/openai_text" >}})
content provider to generate text with ChatGPT through OpenAI API.

### Installation

Before using [`openai_text`]({{< ref "plugins/openai/content-providers/openai_text" >}}) content
provider, it's necessary to add [`blackstork/openai`]({{< ref "plugins/openai" >}}) plugin as a
dependency and install it locally.

To start, update the `hello.fabric` file with the global configuration block:

```hcl
fabric {
  plugin_versions = {
    "blackstork/openai" = ">= 0.4.0"
  }
}

document "greeting" {

  vars {
    solar_system = {
      planets = [
        "Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune"
      ]
      moons_count = 146
    }
  }

  title = "The Greeting"

  content text {
    value = "Hello, Fabric!"
  }

  content text {
    local_var = query_jq(".vars.solar_system.planets | length")

    value = <<-EOT
      There are {{ .vars.local }} planets and {{ .vars.solar_system.moons_count }} moons in our solar system.
    EOT
  }

}
```

Here, a plugin name `blackstork/openai` was added to the list of dependencies in the `plugin_versions` argument
in [the global configuration]({{< ref "language/configs.md#global-configuration" >}}).

With the `hello.fabric` file updated, you can install all required plugins with the `fabric install` command:

```shell
$ fabric install
Mar 11 19:20:10.769 INF Searching plugin name=blackstork/openai constraints=">=v0.4.0"
Mar 11 19:20:10.787 INF Installing plugin name=blackstork/openai version=0.4.0
$
```

Fabric fetched the `blackstork/openai` plugin release from the plugin registry and installed it in the local `./.fabric/` folder.

### Configuration

[`openai_text`]({{< ref "plugins/openai/content-providers/openai_text" >}}) content provider requires an OpenAI API key. The key can be set in the provider's configuration block. It's recommended to store credentials and API keys separately from Fabric code, and using the `env` object to read the key from the `OPENAI_API_KEY` environment variable.

The `config` block for the `openai_text` content provider would look like this:

```hcl
config content openai_text {
  api_key = env.OPENAI_API_KEY
}
```

Add this block to `hello.fabric` file.

### Usage

Lets define the content block that uses `openai_text` content provider:

```hcl
// ...

document "greeting" {

  // ...

  content openai_text {
    local_var = query_jq("{planet: .vars.solar_system.planets[-1]}")

    prompt = <<-EOT
      Share a fact about the planet specified in the provided data:
      {{ .vars.local | toRawJson }}
    EOT
  }

}
```

A JQ query `"{planet: .vars.solar_system.planets[-1]}` fetches the last item from the list
(`Neptune`) and creates a new JSON object `{"planet": "Neptune"}`. The results of the query
execution are stored in `local` variable in the context.

{{< hint note >}}
If you would like to specify a system prompt for OpenAI API, you can set it up in the configuration for `openai_text` provider. See the provider's [documentation]({{< ref "plugins/openai/content-providers/openai_text" >}}) for more configuration options.
{{< /hint >}}

The complete content of the `hello.fabric` file should look like this:

FIXME: TORUN
```hcl
fabric {
  plugin_versions = {
    "blackstork/openai" = ">= 0.4"
  }
}

config content openai_text {
  api_key = env.OPENAI_API_KEY
}

document "greeting" {

  vars "solar_system" {
    planets = [
      "Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune"
    ]

    moons_count = 146
  }

  title = "The Greeting"

  content text {
    value = "Hello, Fabric!"
  }

  content text {
    local_var = query_jq(".vars.solar_system.planets | length")

    value = <<-EOT
      There are {{ .vars.local }} planets and {{ .vars.solar_system.moons_count }} moons in our solar system.
    EOT
  }

  content openai_text {
    local_var = query_jq("{planet: .vars.solar_system.planets[-1]}")

    prompt = <<-EOT
      Share a fact about the planet specified in the provided data:
      {{ .vars.local | toRawJson }}
    EOT
  }

}
```

To render the document, `OPENAI_API_KEY` environment variable must be provided. We can set it for `fabric` command execution:

```shell
$ OPENAI_API_KEY="<key-value>" fabric render document.greeting
...
```

{{< hint warning >}}
Remember to replace `<key-value>` in the CLI command with your OpenAI API key value.
{{< /hint >}}

The results of the render should look similar to the following:

```bash
$ OPENAI_API_KEY="<key-value>" ./fabric render document.greeting
Mar 11 20:39:17.834 INF Loading plugin name=blackstork/openai path=.fabric/plugins/blackstork/openai@0.4.0
# The Greeting

Hello, Fabric!

There are 8 planets and 146 moons in our solar system.

Neptune is the eighth planet from the Sun in our solar system and is the coldest planet. It has average temperatures of minus 353 degrees Fahrenheit (minus 214 degrees Celsius).
```

## Publishing

We can already use the Markdown output (and render it with a Markdown editor like [MacDown](https://macdown.uranusjr.com/)), but it might be better to produce formatted HTML and PDF documents.

To format the document as HTML or PDF, and publish them to a local or an external destination, we use `publish` blocks.

Add two `publish` blocks to the document template:

```hcl
document "greeting" {

  // ...

  // Creating publishing to a local PDF file
  publish local_file {
    path = "./greeting-{{ now | date \"2006_01_02\" }}.{{.format}}"
    format = "pdf"
  }

  // Creating publishing to a local HTML file
  publish local_file {
    path = "./greeting-{{ now | date \"2006_01_02\" }}.{{.format}}"
    format = "html"
  }

}
```

Note, that both `publish` blocks use Go template strings as values for the `path` arguments. This
allows us to have current date in a filename of the output files.

# Next steps

Congratulations! By completing this tutorial, you've gained a good understanding of Fabric and its core principles.

Take a look at the detailed [FCL specification]({{< ref "language" >}}), explore [the open-source templates]({{< ref "templates" >}}) the community made, and see if there are integrations for your tech stack in [Fabric plugins]({{< ref "plugins" >}}).

If you have any questions, feel free to ask in the [Fabric Community Slack](https://fabric-community.slack.com/) and we'll be glad to assist you!
