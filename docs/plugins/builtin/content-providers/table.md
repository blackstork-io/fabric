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

Each cell template has access to the data context and the following variables:
* `.block.rows` – the value of rows_var attribute
* `.block.row` – the current row from `.block.rows` list
* `.block.row_index` – the current row index
* `.block.col_index` – the current column index

Header templates have access to the same variables as value templates,
except for `.block.row` and `.block.row_index`

The content provider is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.


#### Configuration

The content provider doesn't support any configuration parameters.

#### Usage

The content provider supports the following execution parameters:

```hcl
content table {
  # A list of objects representing rows in the table.
  # May be set statically or as a result of one or more queries.
  #
  # Required data.
  # Must have a length of at least 1
  # For example:
  rows_var = null

  # List of header and value go templates for each column
  #
  # Required list of object.
  # Must have a length of at least 1
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

