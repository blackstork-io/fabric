---
title: "`table` content provider"
plugin:
  name: blackstork/builtin
  description: "Produces a table"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: content-provider
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.2" "table" "content provider" >}}

## Description
Produces a table.

Each cell template has access to the data context and the following variables:
* `.rows` – the value of `rows` argument
* `.row.value` – the current row from `.rows` list
* `.row.index` – the current row index
* `.col.index` – the current column index

Header templates have access to the same variables as value templates,
except for `.row.value` and `.row.index`

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration arguments.

#### Usage

The content provider supports the following execution arguments:

```hcl
content table {
  # A list of objects representing rows in the table.
  # May be set statically or as a result of one or more queries.
  #
  # Optional data.
  # Default value:
  rows = null

  # List of header and value go templates for each column
  #
  # Required list of object.
  # Must be non-empty
  # For example:
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

