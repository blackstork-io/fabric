---
title: "`json` data source"
plugin:
  name: blackstork/builtin
  description: "Loads JSON files with the names that match a provided `glob` pattern or a single file from a provided `path`value"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.2" "json" "data source" >}}

## Description
Loads JSON files with the names that match a provided `glob` pattern or a single file from a provided `path`value.

Either `glob` or `path` argument must be set.

When `path` argument is specified, the data source returns only the content of a file.
When `glob` argument is specified, the data source returns a list of dicts that contain the content of a file and file's metadata. For example:

```json
[
  {
    "file_path": "path/file-a.json",
    "file_name": "file-a.json",
    "content": {
      "foo": "bar"
    }
  },
  {
    "file_path": "path/file-b.json",
    "file_name": "file-b.json",
    "content": [
      {"x": "y"}
    ]
  }
]
```

The data source is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

## Configuration

The data source doesn't support any configuration arguments.

## Usage

The data source supports the following execution arguments:

```hcl
data json {
  # A glob pattern to select JSON files to read
  #
  # Optional string.
  # For example:
  # glob = "path/to/file*.json"
  # 
  # Default value:
  glob = null

  # A file path to a JSON file to read
  #
  # Optional string.
  # For example:
  # path = "path/to/file.json"
  # 
  # Default value:
  path = null
}
```