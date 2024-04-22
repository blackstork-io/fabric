---
title: table
plugin:
  name: blackstork/builtin
  description: "Produces a table"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "table" "content provider" >}}

## Description
Produces a table.

This content provider assumes that `query_result` is a list of objects representing rows,
and uses the configured `value` go templates (see below) to display each row.

NOTE: `header` templates are executed with the whole context availible, while `value`
templates are executed on each item of the `query_result` list.

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration parameters.

#### Usage

The content provider supports the following execution parameters:

```hcl
content table {
  # List of header and value go templates for each column
  #
  # Required list of object. For example:
  columns = [{
    header = "1st column header template"
    value  = "1st column values template"
    }, {
    header = "2nd column header template"
    value  = "2nd column values template"
    }, {
    header = "..."
    value  = "..."
  }]
}
```

