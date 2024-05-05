---
title: json
plugin:
  name: blackstork/builtin
  description: "Loads JSON files with the names that match a provided \"glob\" pattern or a single file from a provided path"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "json" "data source" >}}

## Description
Loads JSON files with the names that match a provided "glob" pattern or a single file from a provided path.

Either "glob" or "path" attribute must be provided.

When "path" is specified, only the content of the file is returned.
When "glob" is specified, the structure returned by the data source is a list of dicts that contain the content of the file and file metadata.

For example, with "glob" set, the data source will return the following data structure:

```json
[
  {
    "file_path": "path/file-a.json",
    "file_name": "file-a.json",
    "content": {
      "foo": "bar"
    },
  },
  {
    "file_path": "path/file-b.json",
    "file_name": "file-b.json",
    "content": [
      {"x": "y"}
    ],
  }
]
```

The data source is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

## Configuration

The data source doesn't support configuration.

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data json {
  # A glob pattern to select JSON files for reading
  #
  # For example:
  # glob = "data/*_alerts.json"
  #
  # Optional string. Default value:
  glob = null

  # A file path to a JSON file to read
  #
  # For example:
  # path = "data/alerts.json"
  #
  # Optional string. Default value:
  path = null
}
```