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

Throughout this tutorial, the command examples were executed in macOS Sonoma, in `zsh` shell.

## Hello, Fabric

Let's start with a straightforward "Hello, Fabric!" template to confirm that everything is configured correctly.

Create a new `hello.fabric` file and define a simple template:

```hcl
document "greeting" {

  content text {
    value = "Hello, Fabric!"
  }

}
```

In this code snippet, the `document.greeting` block defines a template with a single anonymous content block containing the text "Hello, Fabric!"

To render the document, ensure that Fabric can locate the `hello.fabric` file. Execute `fabric` command in the same directory (or explicitly provide `--source-dir` CLI argument):

```shell
fabric render document.greeting
```

The output should resemble the following:

```shell
$ fabric render document.greeting
Hello, Fabric!
```

## Document title

Documents typically include titles, and the document block supports the `title` argument as a straightforward method to set a title for a document.

Enhance the `document.greeting` template by adding a title using the `title` argument:

```hcl
document "greeting" {

  title = "The Greeting"

  content text {
    value = "Hello, Fabric!"
  }

}
```

{{< hint note >}}
`title` argument for `document` block is a syntactic sugar that Fabric translates into `content.title` block:

```hcl
content title {
  value = "The Greeting"
}
```

See [the documentation]({{< ref "plugins/builtin/content-providers/text" >}}) for the details about the arguments supported by `text` content provider.
{{< /hint >}}

The rendered output should now include the document title:

```markdown
$ fabric render document.greeting
# The Greeting

Hello, Fabric!
```

## Data blocks

A core feature of the Fabric configuration language is the ability to define data requirements inside templates with the [data blocks]({{< ref "language/data-blocks.md" >}}). The easiest way is to use [`inline`]({{< ref "plugins/builtin/data-sources/inline" >}}) data source that supports free-form data structures.

Note, you must define `data` blocks on the root level of the `document` block.

Modify the template in the `hello.fabric` file to include a `data` block and another `content.text` block:

```hcl
document "greeting" {

  data inline "solar_system" {
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
    query = ".data.inline.solar_system.planets | length"

    value = <<-EOT
      There are {{ .query_result }} planets and {{ .data.inline.solar_system.moons_count }} moons in our solar system.
    EOT
  }

}
```

The content blocks can access and transform the data available with [JQ query](https://jqlang.github.io/jq/manual/) in the [`query`]({{< ref "language/content-blocks.md#generic-arguments" >}}) argument. It's applied to [the context object]({{< ref "language/content-blocks.md#context" >}}) during the evaluation of the block, and the result is stored in the context object under the `query_result` field.

As you can see, `text` argument value in the new content block is a template string – `content.text` blocks support [Go templates](https://pkg.go.dev/text/template) out-of-the-box. The templates can access the context object, so it's easy to include `query_result` or `moons_count` values from the context.

The rendered output should now include the new sentence:

```shell
$ fabric render document.greeting
# The Greeting

Hello, Fabric!

There are 8 planets and 146 moons in our solar system.
```

## Content providers

Fabric seamlessly integrates with external APIs for content generation. An excellent example is the use of the OpenAI API to dynamically generate text through prompts.

In scenarios where providing the exact text or a template string for the content block proves challenging or impossible, leveraging generative AI for summarization becomes invaluable. This enables users to dynamically create context-aware text.

In this tutorial, we will use the [`openai_text`]({{< ref "plugins/openai/content-providers/openai_text" >}}) content provider to generate text with the OpenAI Large Language Model (LLM).

### Installation

Before using [`openai_text`]({{< ref "plugins/openai/content-providers/openai_text" >}}) content provider, it's necessary to add [`blackstork/openai`]({{< ref "plugins/openai" >}}) plugin as a dependency and install it.

To achieve this, update the `hello.fabric` file to resemble the following:

```hcl
fabric {
  plugin_versions = {
    "blackstork/openai" = ">= 0.4.0"
  }
}

document "greeting" {

  data inline "solar_system" {
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
    query = ".data.inline.solar_system.planets | length"

    value = <<-EOT
      There are {{ .query_result }} planets and {{ .data.inline.solar_system.moons_count }} moons in our solar system.
    EOT
  }

}
```

Here, we added a fully qualified name plugin name `blackstork/openai` to the list of dependencies in the `plugin_versions` argument in [the global configuration]({{< ref "language/configs.md#global-configuration" >}}).

With the `hello.fabric` file updated, you can install all required plugins with the `fabric install` command:

```shell
$ fabric install
Mar 11 19:20:10.769 INF Searching plugin name=blackstork/openai constraints=">=v0.4.0"
Mar 11 19:20:10.787 INF Installing plugin name=blackstork/openai version=0.4.0
$
```

Fabric fetched the `blackstork/openai` plugin release from the plugin registry and installed it in the `./.fabric/` folder.

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
...

document "greeting" {

  ...

  content openai_text {
    query = "{planet: .data.inline.solar_system.planets[-1]}"
    prompt = <<-EOT
      Share a fact about the planet specified in the provided
      data: {{ .query_result | toRawJson }}
    EOT
  }

}
```

A JQ query `"{planet: .data.inline.solar_system.planets[-1]}` in `query` argument fetches the last item from the list (`Neptune`) and creates a new JSON object `{"planet": "Neptune"}`, stored under `query_result` field in the context.

{{< hint note >}}
If you would like to specify a system prompt for OpenAI API, you can set it up in the configuration for `openai_text` provider. See the provider's [documentation]({{< ref "plugins/openai/content-providers/openai_text" >}}) for more configuration options.
{{< /hint >}}

The complete content of the `hello.fabric` file should look like this:

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

  data inline "solar_system" {
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
    query = ".data.inline.solar_system.planets | length"

    value = <<-EOT
      There are {{ .query_result }} planets and {{ .data.inline.solar_system.moons_count }} moons in our solar system.
    EOT
  }

  content openai_text {
    query = "{planet: .data.inline.solar_system.planets[-1]}"
    prompt = <<-EOT
      Share a fact about the planet specified in the provided
      data: {{ .query_result | toRawJson }}
    EOT
  }

}
```

To render the document, `OPENAI_API_KEY` environment variable must be set. A simple way to do that, is to set it for each `fabric` command execution:

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

## Markdown rendering

Fabric produces Markdown documents that are compatible with various Markdown editors, allowing rendering in formats such as HTML or PDF. It's also possible to copy-paste rich text into the word processors like Microsoft Word or Google Docs.

An excellent choice for macOS users is [MacDown](https://macdown.uranusjr.com/), an open-source Markdown editor.

![Rendered template in the MacDown Markdown editor](/images/the-greeting.png "The document in MacDown editor")

# Next steps

Congratulations! By completing this tutorial, you've gained a good understanding of Fabric and its core principles.

Take a look at the detailed [FCL specification]({{< ref "language" >}}), explore [the open-source templates]({{< ref "templates" >}}) the community made, and see if there are integrations for your tech stack in [Fabric plugins]({{< ref "plugins" >}}).

If you have any questions, feel free to ask in the [Fabric Community Slack](https://fabric-community.slack.com/) and we'll be glad to assist you!
