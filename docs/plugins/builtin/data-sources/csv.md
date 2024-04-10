---
title: csv
plugin:
  name: blackstork/builtin
  description: "Imports and parses a csv file"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "csv" "data source" >}}

## Description
Imports and parses a csv file.

We assume the table has a header and turn each line into a map based on the header titles.

For example following table

| column_A | column_B | column_C |
| -------- | -------- | -------- |
| Test     | true     | 42       |
| Line 2   | false    | 4.2      |

will be represented as the following structure:
```json
[
  {"column_A": "Test", "column_B": true, "column_C": 42},
  {"column_A": "Line 2", "column_B": false, "column_C": 4.2}
]
```

The data source is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

## Configuration

The data source supports the following configuration parameters:

```hcl
config "data" "csv" {
  # Must be a one-character string
  #
  # Optional. Default value:
  delimiter = ","
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data "csv" {
  # Required. For example:
  path = "path/to/file.csv"
}
```