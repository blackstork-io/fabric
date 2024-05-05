---
title: csv
plugin:
  name: blackstork/builtin
  description: "Loads CSV files with the names that match a provided \"glob\" pattern or a single file from a provided path"
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
Loads CSV files with the names that match a provided "glob" pattern or a single file from a provided path.

Either "glob" or "path" attribute must be set.

When "path" attribute is specified, only the content of the file is returned. For example, the following CSV data:

| column_A | column_B | column_C |
| -------- | -------- | -------- |
| Foo      | true     | 42       |
| Bar      | false    | 4.2      |

```json
[
  {"column_A": "Foo", "column_B": true, "column_C": 42},
  {"column_A": "Bar", "column_B": false, "column_C": 4.2}
]
```

When "glob" attribute is specified, the structure returned by the data source is a list of dicts that contain the content of the file and file metadata. For example:

```json
[
  {
    "file_path": "path/file-a.csv",
    "file_name": "file-a.csv",
    "content": [
      {"column_A": "Foo", "column_B": true, "column_C": 42},
      {"column_A": "Bar", "column_B": false, "column_C": 4.2}
    ]
  },
  {
    "file_path": "path/file-b.csv",
    "file_name": "file-b.csv",
    "content": [
      {"column_C": "Baz", "column_D": 1},
      {"column_C": "Clu", "column_D": 2}
    ]
  },
]
```

The data source is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data csv {
  # Must be a one-character string
  #
  # Optional string. Default value:
  delimiter = ","
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data csv {
  # A glob pattern to select CSV files for reading
  #
  # For example:
  # glob = "path/to/files*.csv"
  #
  # Optional string. Default value:
  glob = null

  # A file path to a CSV file to read
  #
  # For example:
  # path = "path/table.csv"
  #
  # Optional string. Default value:
  path = null
}
```