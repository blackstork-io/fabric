---
title: Content Blocks
description: Learn how to use Fabric content blocks efficiently for building modular and reusable document templates.
type: docs
weight: 60
---
# Content blocks

`content` blocks define document segments: text paragraphs, tables, graphs, lists, etc.

The block signature includes the name of the content provider that will execute the content block.

```hcl
# Root-level definition of a content block
content <content-provider-name> "<block-name>" {
  # ...
}

document "foobar" {

  # In-document named definition of a content block
  content <content-provider-name> "<block-name>" {
    # ...
  }

  # In-document anonymous definition of a content block
  content <content-provider-name> {
    # ...
  }

}
```

The order of the `content` blocks in the template determines the order of the generated content in
the document.

If the block is at the root level of the file, outside of the `document` block, both names –
the content provider name and the block name – are required. A combination of block type `content`,
content provider name, and block name serves as a unique identifier of a block within the codebase.

If the content block is defined within the document template, only a content provider name is
required and a block name is optional.

A content block is rendered by a specified content provider. See [Content Providers]({{< ref
"content-providers.md" >}}) for the list of the content providers supported by Fabric.

## Supported arguments

The arguments provided in the block are either generic arguments or provider-specific arguments.

### Generic arguments

- `config`: (optional) a reference to a named configuration block for the content provider. If
  provided, it takes precedence over the default configuration. See content provider
  [configuration details]({{< ref "configs.md#block-configuration" >}}) for more information.
- `local_var`: (optional) a shortcut for specifying a local variable. See [Variables]({{< ref
  "context.md#variables" >}}) for the details.

### Content provider arguments

Content provider arguments differ per content provider. See the documentation for a specific content provider (find it in [Content Providers]({{< ref "content-providers.md" >}}) documentation) for the details on the supported arguments.

## Supported nested blocks

- `meta`: (optional) a block containing metadata for the block. See [Metadata]({{< ref "configs.md#metadata" >}}) for details.
- `config`: (optional) an inline configuration for the block. If provided, it takes precedence over the `config` argument and default configuration for the content provider.
- `vars`: (optional) a block with variable definitions. See [Variables]({{< ref
  "context.md#variables" >}}) for the details.

## References

See [References]({{< ref references.md >}}) for the details about referencing content blocks.

## Example

```hcl
config content openai_text "test_account" {
  # Reading a key from an environment variable
  api_key = env.OPENAI_API_KEY
}

document "test-doc" {

  vars {
    items = ["aaa", "bbb", "ccc"]
  }

  content text {
    # Query contains a JQ query executed against the context
    local_var = query_jq(".vars.items | length")

    # The context can be accessed in Go templates
    value = "There are {{ .vars.local }} items: {{ .vars.items | toPrettyJson }}"
  }

  content openai_text {
    config = config.content.openai_text.test_account

    prompt = <<-EOT
       Write a short story, just a paragraph, about space exploration
       using the values from the provided items list as character names:

       {{ .vars.items | toPrettyJson }}
    EOT
  }
}
```

produces the following output:

```text
There are 3 items: [
  "aaa",
  "bbb",
  "ccc"
]

In the vast expanse of outer space, aaa, bbb, and ccc embarked on a daring mission of exploration.
Their spaceship soared through the galaxies, encountering unknown planets and celestial bodies. With
aaa's courage leading the way, bbb's wisdom guiding their decisions, and ccc's ingenuity solving the
complex challenges they faced, the trio delved deeper into the mysteries of the cosmos. Together,
they pushed the boundaries of what was thought possible, and with each discovery, their bond grew
stronger. In the endless sea of stars, they found not only the wonders of the universe but also the
strength of their friendship.
```

## Next steps

See [Section Blocks]({{< ref "section-blocks.md" >}}) documentation to learn how to group the content into the sections for better maintainability and re-usability.
