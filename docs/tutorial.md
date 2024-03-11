---
title: Tutorial
type: docs
weight: 30
---

# Tutorial

This tutorial will give you everything you need to use Fabric and [Fabric Configuration Language]({{< ref "language" >}}) (FCL) to generate documents. We'll walk through creating a simple template, adding data blocks, filtering and mutating data, installing plugins, and rendering text using external providers.

## Prerequisites

To follow the tutorial, you will need:

- Fabric CLI [installed]({{< ref "install.md" >}}) and `fabric` CLI command available
- (optional) OpenAI API token

Throughout this tutorial, the command examples were executed in macOS Sonoma, in `zsh` shell.

## Hello, Fabric

Let's start with a simple "Hello, Fabric!" template to make sure everything is set up correctly.

Create a new `hello.fabric` file and define a simple template:

```hcl
document "greeting" {

  content text {
    text = "Hello, Fabric!"
  }

}
```

Here, `document.greeting` block defines a template with a single anonymous content block with "Hello, Fabric!" text.

To render the document, make sure Fabric can find `hello.fabric` file, by executing `fabric` in the same directory or providing `–source-dir`, and run:

```shell
fabric render document.greeting
```

If everything was installed correctly, the output should look like this:

```shell
$ fabric render document.greeting
Hello, Fabric!
```

## Document title

The documents usually have titles. `document` block supports `title` argument - a simple way to set a title for a document.

Add a title to `document.greeting` template using `title` argument:

```shell
document "greeting" {

  title = "The Greeting"

  content text {
    text = "Hello, Fabric!"
  }

}
```

{{< hint note >}}
`title` argument for `document` block is a syntactic sugar that Fabric translates into `content.text` block with `format_as` attribute set to `title`:

```hcl
content text {
  text = "The Greeting"
  format_as = "title"
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

A core feature of Fabric configuration language is the ability to define data requirements inside the templates with the [data blocks]({{< ref "language/data-blocks.md" >}}). The easiest way is to use [`inline`]({{< ref "plugins/builtin/data-sources/inline" >}}) data source that supports free-form data structures.

Note, you must define `data` blocks on the root level of `document` block.

Change the template in `hello.fabric` file to include a `data` block and another `content.text` block:

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
    text = "Hello, Fabric!"
  }

  content text {
    query = ".data.inline.solar_system.planets | length"

    text = <<-EOT
      There are {{ .query_result }} planets and {{ .data.inline.solar_system.moons_count }} moons
      in our solar system.
    EOT
  }

}
```

The content blocks can access and transform the data available with [JQ query](https://jqlang.github.io/jq/manual/) in the [`query`]({{< ref "language/content-blocks.md#generic-arguments" >}}) argument. It's applied to [the context object]({{< ref "language/content-blocks.md#context" >}}) during the evaluation of the block, and the result is stored in the context object under `query_result` field.

As you can see, `text` argument value in the new content block is a template string – `content.text` blocks support [Go templates](https://pkg.go.dev/text/template) out-of-the-box. The templates can access the context object, so it's easy to include `query_result` or `moons_count` values from the context.

The rendered output should now include the new sentence:

```shell

$ fabric render document.greeting
# The Greeting

Hello, Fabric!

There are 8 planets and 146 moons in our solar system.
```

## External content providers

Fabric also integrates with external APIs for content generation. For example, it's possible to use the OpenAI API to generate text dynamically with prompts.

In certain scenarios, providing the exact text or a template string for the content block might be difficult or impossible. In such cases, leveraging generative AI for summarization proves beneficial, enabling users to create context-aware text dynamically.

### Installation

Before we can use it, we must add `blackstork/openai` plugin as a dependency and install it.

Update the list of plugin dependencies by adding a fully qualified name `blackstork/openai` for the [OpenAI plugin]({{< ref "plugins/openai" >}}) to `plugin_versions` dict in the [global configuration block]({{< ref "language/configs.md#global-configuration" >}}). The plugin also requires a [version constraint](https://semver.org/).

Add the following code to `hello.fabric` file:

```hcl
fabric {
  plugin_versions = {
    "blackstork/openai" = ">= 0.4"
  }
}
```

With `hello.fabric` updated, it's easy to install all required plugins with `fabric install` command:

```shell
$ fabric install

TBD
```

If the command succeeded, Fabric fetched `blackstork/openai` plugin release from the Fabric plugin registry and installed it in `./.fabric/plugins` folder.

### Configuration

`blackstork/openai` plugin contains `openai_text` content provider that uses OpenAI [Chat Completions API](https://platform.openai.com/docs/guides/text-generation/chat-completions-api) for text generation.

`openai_text` content provider requires an OpenAI API key, which can be set in the provider's configuration block. It's a good practice to store the credentials and API keys separate from the code, so we can use `from_env_variable` function to read the key from environment variable `OPENAI_API_KEY`.

`config` block for `openai_text` content provider would look like this:

```hcl
config content openai_text {
  api_key = from_env_variable("OPENAI_API_KEY")
}
```

### Usage

Now we can define a simple content block powered by `openai_text` content provider (see the full version of the code below):

```hcl
...

document "greeting" {

  ...

  content openai_text {
    query = "{planet: .data.inline.solar_system.planets[-1]}"
    prompt = "Share a fact about the planet specified in the provided data"
  }

}
```

Here, a JQ query `"{planet: .data.inline.solar_system.planets[-1]}` performs two operations: it fetches the last item in the list (`Neptune`) and returns a new JSON object `{"planet": "Neptune"}`. This object is stored under `query_result` field in the context. Combined with the `prompt` string, `query_result` value creates a user prompt for OpenAI API.

{{< hint note >}}
If you would like to specify a system prompt for OpenAI API, you can set it up in the configuration for `openai_text` provider. See the provider's [documentation]({{< ref "plugins/openai/content-providers/openai_text" >}}) for more configuration options.
{{< /hint >}}

The full content of `hello.fabric` file should look like this:

```hcl
fabric {
  plugin_versions = {
    "blackstork/openai" = ">= 0.4"
  }
}

config content openai_text {
  api_key = from_env_variable("OPENAI_API_KEY")
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
    text = "Hello, Fabric!"
  }

  content text {
    query = ".data.inline.solar_system.planets | length"

    text = <<-EOT
      There are {{ .query_result }} planets and {{ .data.inline.solar_system.moons_count }} moons
      in our solar system.
    EOT
  }

  content openai_text {
    query = "{planet: .data.inline.solar_system.planets[-1]}"
    prompt = "Share a fact about the planet specified in the provided data"
  }

}
```

To render the document, `OPENAI_API_KEY` environment variable must be set. A simple way to do that, is to set it for per execution:

```shell
$ OPENAI_API_KEY="<key-value>" fabric render document.greeting
...
```

{{< hint warning >}}
Remember to replace `<key-value>` in the CLI command with your OpenAI API key value.
{{< /hint >}}

# What's next

Congratulations! By completing this tutorial, you've gained a solid understanding of Fabric and {{< dfn "FCL" >}} core principles.

Take a look at the open-source templates the community built in [Fabric Templates](https://github.com/blackstork-io/fabric-templates) GitHub repository – you can reuse the whole documents or specific blocks in your own templates!

If you have any questions, share them in the [Fabric Community Slack](https://fabric-community.slack.com/) and we will be glad to help!
