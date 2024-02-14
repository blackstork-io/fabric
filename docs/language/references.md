---
title: References
type: docs
weight: 80
---

# References

Fabric language supports code reuse via block references. Note, only `data`, `content` and `section` blocks with names, defined on a root level of the file, can be referenced.

To include a named block defined on the root level into a document, use `ref` label and `base` argument:

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
    base = <block-identifier>
    ... 
  }
}
```

- `<block-type>` is either `content`, `data` or `section`.
- `<block-identifier>` is a dot-separated identifier for the block to be included. It consists of all labels in the block signature. For example, `content.<content-provider>.<block-name>` or `section.<block-name>`. The block name in the identifier shouldn't be wrapper in double quotes `"`.

If the `ref` block is defined on the root level of the file, the block name is required. If it's within the `document` block, the block can be anonymous - the name is optional.

{{< hint warning >}}
An anonymous `ref` block adopts a name of the referenced block. Since the final name of the block isn't stated explicitly, unwanted overrides can happen — a block signature must be unique between the blocks defined on the same level.
{{< /hint >}}

## Overriding arguments

The `ref` block definition can include the attribute that would override the attributes provided in the original block. This is helpful if the block's behaviour needs adjustments per document.

## Query input requirement

Content blocks rely on `query` attribute for selecting data needed for rendering (see content blocks' [Generic Arguments]({{< ref "content-blocks.md#generic-arguments" >}})). The JQ query uses the data path which is often document-specific and depends on the name of the data block. This hinders the re-usability of the content blocks.

Fabric supports an explicit way for the content block to require the input data - `query_input` and `query_input_required` attributes. If `query_input_required` set to `true`, the content block expects `query_input` attribute to be provided in the `ref` block.

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

  text = "The data contains {{ .query_result }} elements"
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
