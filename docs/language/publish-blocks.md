---
title: Publish Blocks
description: Learn how to use Fabric publish blocks to define the destinations for produced documents.
type: docs
weight: 72
---

# Publish blocks

## Overview

`publish` blocks define the format and the destinations for document delivery.

After rendering, the document is published to a local (for example, a file on a filesystem) or an
external destination (for example, Google Drive or GitHub), formatted as Markdown, PDF, or HTML.
`publish` blocks, similar to `data` blocks, define the integrations. In this case the integrations
are responsible for document delivery.

Similar to `data` and `content` blocks, `publish` block signature includes the name of the publisher
that will execute the integration:

```hcl
document "foobar" {

  # In-document named definition of a publish block
  publish <publisher-name> "<block-name>" {
    # ...
  }

  # In-document anonymous definition of a publish block
  publish <publisher-name> {
    # ...
  }

}
```

If `publish` block is placed at the root level of the file, outside of the `document` block, both
names – the publisher name and the block name – are required. A combination of block type `publish`,
a publisher name, and a block name serves as a unique identifier of a block within the codebase.

If `publish` block is defined within the document, only a publisher name is needed and a block name
is optional.

Every `publish` block is executed by a corresponding publisher. See [Publishers]({{< ref
publishers.md >}}) for the list of supported publishers.

Fabric executes `publish` blocks as the last step of the processing, after the document is rendered
and ready for formatting and delivery.

## Formatting

Fabric supports a set of formatting options for the output documents: Markdown, PDF, and HTML.

The publishers declare the formats they support (see the documentation for a specific publisher
([Publishers]({{< ref publishers.md >}}) for more information). For example, [`local_file`]({{< ref
"local_file.md" >}}) publisher supports all three format types: `md`, `pdf` and `html`

### HTML formatting

The template authors can configure HTML formatting: to add JS script and CSS script tags, or include
JS or CSS code inline.

To customize the produced HTML document, use `frontmatter` content block on the root level of the
document template.

The supported fields are:

- `title` — a string, used as a HTML page title. If not set, the formatter will use the first title
  from the template. If the template has no title elements, "Untitled" value will be used.
- `description` — a string, used as a HTML page description meta tag value
- `js_sources` — a list of strings, included in HTML document head as `<script async defer type="application/javascript" src="[value]"></script>` tags
- `css_sources` — a list of strings, included in HTML document head as `<link type="text/css" rel="stylesheet" href="[value]" />` tags
- `js_code` — a string, included in HTML document head as a body of `<script type="text/javascript">[value]</script>` tag
- `css_code` — a string, included in HTML document head as a body of `<style>[value]</style>` tag

All fields are optional.

For example:

```hcl
document "test" {

  content frontmatter {
    content = {
      title = "Foo Title"
      description = "Bar Description"

      js_sources = ["https://buttons.github.io/buttons.js", "/static/local.js"]
      css_sources = ["/static/main.css", "https://localhost.localhost/some.css"]

      js_code = <<-EOT
        console.info("JS code execution");
      EOT

      css_code = <<-EOT
        a {
          font-family: Verdana;
        }
      EOT
    }
    format = "yaml"
  }

  title = "Main Document Title"

  content text {
    value = "Test Body"
  }

  publish local_file {
    path = "./test-document.html"
    format = "html"
  }
}
```

The template, when rendered and published, will produce `./test-document.html` file containing:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="Bar Description">
    <title>Foo Title</title>
    <script async defer type="application/javascript" src="https://buttons.github.io/buttons.js"></script>
    <script async defer type="application/javascript" src="/static/local.js"></script>
    <link type="text/css" rel="stylesheet" href="/static/main.css" />
    <link type="text/css" rel="stylesheet" href="https://localhost.localhost/some.css" />
    <script type="text/javascript">
        console.info("JS code execution");
    </script>
    <style>
     a {
       font-family: Verdana;
     }
    </style>
</head>
<body>
 <h1 id="main-document-title">Main Document Title</h1>
<p>Test Body</p>
</body>
</html>
```

## Supported arguments

The arguments supported in the `publish` block are either generic arguments or publisher-specific
arguments.

### Generic arguments

- `config`: (optional) a reference to a named configuration block for the publisher. If provided, it
  takes precedence over the default configuration. See [Publisher configuration]({{< ref
  "configs.md#publisher-configuration" >}}) for the details.
- `format`: (optional) a format of the output, `md` (Markdown) by default. The publishers declare
  the formats they support. See the documentation for a specific publisher for more information
  ([Publishers]({{< ref publishers.md >}}).

### Publisher arguments

A publisher might define the arguments it supports. See [Publishers]({{< ref publishers.md >}}) for
the details on the supported arguments per publisher.

## Supported nested blocks

- `meta`: (optional) a block containing metadata for the block. See [Metadata]({{< ref
  "configs.md#metadata" >}}) for details.
- `config`: (optional) an inline configuration for the block. If provided, it takes precedence over
  the `config` argument and the default configuration for the publisher.
- `vars`: (optional) a block with variable definitions. See [Variables]({{< ref
  "context.md#variables" >}}) for the details.

## Execution

To execute `publish` blocks, set `--publish` flag when calling `fabric render` command:

```bash
$ fabric render --help

Render the specified document and either publish it or output it as Markdown to stdout.

Usage:
  fabric render TARGET [flags]

Args:
  TARGET   name of the document to be rendered as 'document.<name>'

Flags:
      --format string   default output format of the document (md, html or pdf) (default "md")
  -h, --help            help for render
      --publish         publish the rendered document

Global Flags:
      --color               enables colorizing the logs and diagnostics (if supported by the terminal and log format) (default true)
      --log-format string   format of the logs (plain or json) (default "plain")
      --log-level string    logging level ('debug', 'info', 'warn', 'error') (default "info")
      --source-dir string   a path to a directory with *.fabric files (default ".")
  -v, --verbose             a shortcut to --log-level debug
```

If `render` command called without `--publish` flag, Fabric will render the document and print it to
stdout, either as Markdown or as HTML file, depending on `--format` argument.

For example:

```bash
$ fabric render document.example --format md
# Hello World

Document body

$ fabric render document.example --format html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Hello World</title>
</head>
<body>
<h1 id="document-title">Hell World</h1>
<p>Document body</p>
</body>
</html>%

$ fabric render document.example --publish
Jun  2 12:48:08.899 INF Writing to a file path=/tmp/example_2024_06_02.pdf
Jun  2 12:48:09.182 INF Writing to a file path=/tmp/example_2024_06_02.html
Jun  2 12:48:09.183 INF Writing to a file path=/tmp/example_2024_06_02.md
```

## References

See [References]({{< ref references.md >}}) for the details about referencing `publish` blocks.

## Example

The following document defines two delivery destinations: a local PDF file and a local HTML file.

Note, that [`local_file`]({{< ref "local_file.md" >}}) publisher treats the value of its `path`
argument as a Go template string.

```hcl
document "foo" {

  publish local_file {
    path = "docs/foo_{{ now | date \"2006_01_02\" }}.{{.format}}"
    format = "pdf"
  }

  publish local_file {
    path = "html/foo-latest.{{.format}}"
    format = "html"
  }

  title = "Test Document"

  content text {
    value = "Static text in the document body"
  }
}
```

## Next steps

See [References]({{< ref references.md >}}) for the details about reusing block code via
referencing.
