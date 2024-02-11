---
title: References
type: docs
weight: 80
---

# References

Fabric laguage allows code reuse through references and arguments override. Only a subset of block types supports references: `data`, `content` and `section` block types.

To include a block defined outside the document, on the root level of a file, into a document, use `ref` label and set `base` argument:

```hcl
<block-type> <plugin-name> "<block-name>" {
    ...
}

document "foo" {

    <block-type> ref "<block-name>" {
        base = <block-identifier>

        ... 
    }
}
```

A block identifier is a dot-separated list of block type, plugin name and a block name:
```hcl
<block-type>.<plugin-name>.<block-name>
```

If the `ref` block is defined at the root level of the file, outside of the `document` definition, the block name is required. If the `ref` block is defined within the document, the block name is optional.


{{< hint warning >}}
An anonymous `ref` block adops a name of the referenced block. Since the final name of the block is not stated explicitely, unwanted overrides can happen — a block signature must be unique between the blocks defined on the same level.
{{< /hint >}}


## Argument Overrides

The `ref` block definition can include the attribute that would override the attributes provided in the original block. This is very helpful if the block's behaviour needs to be adjusted per document.


## Query Input Requirement

Content blocks rely on `query` attribute for selecting data needed for rendering (see content blocks' [Generic Arguments]({{< ref "content-blocks.md#generic-arguments" >}})). The JQ query uses the data path which is often document-specific and depends on the name of the data block. This hinders the reusability of the content blocks. 

Fabric supports an explicit way for the content block to require the input data - `query_input` and `query_input_required` attributes. If `query_input_required` is set to `true`, the content block expectes `query_input` attribute to be provided in the `ref` block.


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




