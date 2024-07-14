---
title: Tutorial
description: Dive into Fabric tutorial to learn everything you need to know about using Fabric effectively. From the basics of FCL commands to advanced features such as data configurations, our tutorial provides clear, step-by-step instructions for building Fabric templates. Start improving your workflow with Fabric today.
type: docs
weight: 30
code_blocks_no_wrap: true
---

# Tutorial

This tutorial provides comprehensive guidance on using Fabric and the [Fabric Configuration
Language]({{< ref "language" >}}) (FCL) for document generation. We'll systematically cover creating
a basic template, incorporating data blocks, applying data filtering and mutation, installing
plugins, and rendering text with external content providers.

## Prerequisites

Before you follow this tutorial, make sure you have the following:

- Fabric CLI [installed]({{< ref "install.md" >}}) and `fabric` CLI command available
- (optional) OpenAI API token

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

In this snippet, `document.greeting` block defines a template with a single anonymous content
block with the static text "Hello, Fabric!"

To render the document, execute `fabric` command in the directory with `hello.fabric` file, or
explicitly specify a path to another directory with `--source-dir` CLI argument:

```shell
fabric render document.greeting
```

The command should produce `Hello, Fabric!` string:

```shell
$ fabric render document.greeting
Hello, Fabric!
```

## Document title

Documents usually have titles and the `document` block supports `title` argument as an
easy way to set a title.

With the new `title` argument, `document.greeting` template should look like this:

```hcl
document "greeting" {

  title = "The Greeting"

  content text {
    value = "Hello, Fabric!"
  }

}
```

{{< hint note >}}
`title` argument for `document` block is a syntactic sugar translated into `content.title` block
during rendering:

```hcl
content title {
  value = "The Greeting"
}
```

See [`content.title`]({{< ref "plugins/builtin/content-providers/title" >}}) content provider
documentation for the details.
{{< /hint >}}

The rendered Markdown output should now include the document title:

```markdown
$ fabric render document.greeting
# The Greeting

Hello, Fabric!
```

## Variables

For this tutorial, instead of defining the data requirements with `data` blocks, we will use
variables defined in `vars` block.

Change the template in the `hello.fabric` file to include a `vars` block and add another
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

Here, we defined inline data inside `vars` block (see [Variables]({{< ref
"context.md#variables">}})), used `local_var` (see [Local variable]({{< ref
"context.md#local-variable" >}})), and queried it with `query_jq()` function (see [Querying the
context]({{< ref "context.md#querying-the-context" >}})) inside the `content` block.

The `value` argument in the new content block is a template string – `content.text` blocks support
[Go templates](https://pkg.go.dev/text/template) in `value` argument. The templates can access the
evaluation context, so it's easy to use JSON path and include the values of `local` and
`solar_system.moons_count` variables.

The rendered output will now include the new sentence:

```shell
$ fabric render document.greeting
# The Greeting

Hello, Fabric!

There are 8 planets and 146 moons in our solar system.
```

## Content providers

Fabric uses both internal implementations and integrates with external APIs for content generation.
An excellent example is the use of the OpenAI API to dynamically generate text with prompts.

In scenarios where providing the exact text or a template string for the content block proves
challenging or impossible, we can leverage generative AI for text generation. This allows us to
dynamically create context-aware text.

Lets use [`openai_text`]({{< ref "plugins/openai/content-providers/openai_text" >}}) content
provider for generating text with OpenAI API.

### Installation

Before using [`openai_text`]({{< ref "plugins/openai/content-providers/openai_text" >}}) content
provider, it's necessary to add [`blackstork/openai`]({{< ref "plugins/openai" >}}) plugin as a
dependency and install it locally.

First, update the `hello.fabric` file with the global configuration block:

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

Here, plugin `blackstork/openai` is in the list of dependencies, listed in `plugin_versions` argument
in [the global configuration]({{< ref "language/configs.md#global-configuration" >}}).

With updated `hello.fabric` file, install all required plugins with the `fabric install` command:

```shell
$ fabric install
Mar 11 19:20:10.769 INF Searching plugin name=blackstork/openai constraints=">=v0.4.0"
Mar 11 19:20:10.787 INF Installing plugin name=blackstork/openai version=0.4.0
$
```

Fabric fetched the `blackstork/openai` plugin release from the plugin registry and installed it in
the local `./.fabric/` folder.

The versions in the command output on your system might be different. With `>= 0.4.0` version
constraint, Fabric will install the latest stable version (higher than `0.4.0`) of the plugin.

### Configuration

OpenAI API requires API key for authentication. The key must be set in [`openai_text`]({{< ref
"plugins/openai/content-providers/openai_text" >}}) provider's configuration block. It's recommended
to store credentials separately from Fabric code, and use the `env` object (see [Environment variables]({{< ref
"configs.md#environment-variables" >}})).

We can specify OpenAI API key in `OPENAI_API_KEY` environment variable and access it in Fabric file
with `env.OPENAI_API_KEY`.

The `config` block for the `openai_text` content provider looks like this:

```hcl
config content openai_text {
  api_key = env.OPENAI_API_KEY
}
```

Add this block to the root level of `hello.fabric` file, outside `document` block.

### Usage

Lets define the content block that uses `openai_text` content provider:

```hcl
# ...

document "greeting" {

  # ...

  content openai_text {
    local_var = query_jq("{planet: .vars.solar_system.planets[-1]}")

    prompt = <<-EOT
      Share a fact about the planet specified in the provided data:
      {{ .vars.local | toRawJson }}
    EOT
  }

}
```

In this block we again use `local_var`. A JQ query `"{planet: .vars.solar_system.planets[-1]}`
fetches the last item from the list of planets (`Neptune`) and creates a new JSON object `{"planet":
"Neptune"}`.

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

  content openai_text {
    local_var = query_jq("{planet: .vars.solar_system.planets[-1]}")

    prompt = <<-EOT
      Share a fact about the planet specified in the provided data:
      {{ .vars.local | toRawJson }}
    EOT
  }

}
```

To render the document, set `OPENAI_API_KEY` environment variable when running `fabric` command:

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
fabric render document.greeting
Jun 23 17:14:23.910 INF Parsing fabric files command=render
Jun 23 17:14:23.912 INF Loading plugin resolver command=render includeRemote=false
Jun 23 17:14:23.912 INF Loading plugin runner command=render
Jun 23 17:14:23.939 INF Rendering content command=render target=greeting
Jun 23 17:14:23.939 INF Loading document command=render target=greeting
# The Greeting

Hello, Fabric!

There are 8 planets and 146 moons in our solar system.

Neptune is the eighth and farthest known planet from the Sun in the Solar System. It is classified as an ice giant and is the fourth-largest planet by diameter.
```

## Publishing

By default, Fabric prints rendered document into standard output formatted as Markdown. We can use
any Markdown editor, for example [MacDown](https://macdown.uranusjr.com/) for macOS, to render
Markdown, but it's also possible to produce HTML and PDF documents with Fabric.

To format the document as HTML, PDF, or Markdown, and publish it to a local or an external
destination, use `publish` blocks.

Add a `publish` block to the document template:

```hcl
document "greeting" {

  # ...

  # Publishing to a local HTML file
  publish local_file {
    path = "./greeting-{{ now | date \"2006_01_02\" }}.{{.format}}"
    format = "html"
  }

}
```

Note that similarly to `content.text` blocks, `publish` block supports Go template string as the
`path` argument value. This means we can specify a dynamic path for the output file - in this case,
the filename will contain a date and the output format.

To render the document and publish the output to a local file, use `--publish` flag when running `fabric render`:

```bash
$ fabric render document.greeting --publish
Jun 23 17:28:03.027 INF Parsing fabric files command=render
Jun 23 17:28:03.028 INF Loading plugin resolver command=render includeRemote=false
Jun 23 17:28:03.028 INF Loading plugin runner command=render
Jun 23 17:28:03.056 INF Publishing document command=render target=greeting
Jun 23 17:28:03.056 INF Loading document command=render target=greeting
Jun 23 17:28:04.213 INF Writing to a file command=render path=/tmp/greeting-2024_06_23.html
$
```

You can find the produced HTML file in the current directory (or at the path specified in `path`
argument). The file contents should look similar to this:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>The Greeting</title>
</head>
<body>
 <h1 id="the-greeting">The Greeting</h1>
<p>Hello, Fabric!</p>
<p>There are 8 planets and 146 moons in our solar system.</p>
<p>Neptune is the eighth and most distant planet in our solar system, located about 4.5 billion kilometers away from the Sun.</p>

</body>
</html>
```

To learn hot to add JS and CSS to the produced HTML document, see [Formatting]({{< ref
"publish-blocks.md#formatting" >}}) documentation.

# Next steps

Congratulations! By completing this tutorial, you've gained a good understanding of Fabric and its
core principles.

Take a look at the detailed [FCL specification]({{< ref "language" >}}), explore [the open-source
templates]({{< ref "templates" >}}) the community made, and see if there are integrations for your
tech stack in [Fabric plugins]({{< ref "plugins" >}}).

If you have any questions, feel free to ask in the [Community
Slack](https://fabric-community.slack.com/) and we'll be glad to assist you!
