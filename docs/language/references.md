---
title: References
type: docs
weight: 75
---

# References

Fabric language supports code reuse via block references. It's possible to reference named
`document`, `data`, `content`, `section`, and `publish` blocks defined on a root level of the file.

To include a named block defined on a root level into a document, use `ref` label and `base`
argument:

```hcl
content <content-provider> "<block-name>" {
  ...
}

data <data-source> "<block-name>" {
  ...
}

section "<block-name>" {
  ...
}

document "foo" {

  <block-type> ref "<block-name>" {
    base = <block-identifier-with-matching-block-type>
    ... 
  }

  content ref {
    base = content.<content-provider>.<block-name>
  }

}
```

- `<block-type>` is `document`, `data`, `content`, `section` or `publish`.
- `<block-identifier>` is an identifier for the referenced block. It consists of a dot-separated
  list of all labels in a block signature: `content.<content-provider>.<block-name>` or
  `section.<block-name>`. The block name in the identifier must *not* be wrapper in double quotes `"`.

For the `ref` blocks defined on a root level of the file, the block name is always required. If the
blocks are inside the `document` block, they can be anonymous.

{{< hint warning >}} An anonymous `ref` block adopts a name of the referenced block. Since the final
name of the block isn't stated explicitly, unwanted overrides can happen. A block signature must be
unique between the blocks defined on the same level of the template. {{< /hint >}}

## Overriding arguments

The `ref` block definition can include the argument that would override the arguments provided in
the original block. This is helpful if the block's behaviour needs adjustments per document.

For example:

```hcl
content text "hello_world" {
  value = "Hello, World!"
}

document "foo" {

  content ref "hello_john" {
    base = content.text.hello_world
    value = "Hello, John!"
  }

  content ref {
    base = content.text.hello_world
    value = "Hello, New World!"
  }

}
```

<!-- FIXME: https://github.com/blackstork-io/fabric/issues/29

## Query input requirement

Content blocks rely on `query` argument for selecting data needed for rendering (see content blocks' [Generic Arguments]({{< ref "content-blocks.md#generic-arguments" >}})). The JQ query uses the data path which is often document-specific and depends on the name of the data block. This hinders the re-usability of the content blocks.

Fabric supports an explicit way for the content block to require the input data - `query_input` and `query_input_required` arguments. If `query_input_required` set to `true`, the content block expects `query_input` argument to be provided in the `ref` block.

## Example

```hcl
data elasticsearch "foo" {
  index = "test-index"
  ...
}

content text "qux" {
  # Using `query_input` field in the context that contains the result of
  # the `query_input` query
  query = ".query_input | length"

  # Require the referrer to specify `query_input` query that will be used
  # to get the data for `query_input` field in the context
  query_input_required = true
  value = "The data contains {{ .query_result }} elements"
}

document "test-document" {

  # Anonymous referrer block adops the name of the referenced block - `data.elasticsearch.foo`
  data ref {
    base = data.elasticsearch.foo
  }

  # Named referrer block keeps its name - `data.elasticsearch.bar`
  data ref "bar" {
    base = data.elasticsearch.foo
  }

  # Provided argument `index` overrides the value set in the original block.
  data ref "baz" {
    base = data.elasticsearch.foo
    index = "another-test-index"
  }

  # Referred block requires `query_input` to be provided,
  # so it can be used in query set in `query` argument in the original block.
  content ref {
    base = content.text.qux
    query_input = ".data.elasticsearch.bar"
  }

}
```
-->

## Next steps

To learn how to dynamically adapt the template structure to input data, see [Dynamic Blocks]({{< ref dynamic-blocks.md >}}) documentations.
