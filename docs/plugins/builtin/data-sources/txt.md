---
title: txt
plugin:
  name: blackstork/builtin
  description: "Loads TXT files with the names that match a provided \"glob\" pattern or a single file from a provided path"
  tags: []
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "txt" "data source" >}}

## Description
Loads TXT files with the names that match a provided "glob" pattern or a single file from a provided path.

Either "glob" or "path" attribute must be set.

When "path" is specified, only the content of the file is returned.
When "glob" attribute is specified, the structure returned by the data source is a list of dicts that contain the content of the file and file metadata. For example:

```json
[
  {
    "file_path": "path/file-a.txt",
    "file_name": "file-a.txt",
    "content": "foobar"
  },
  {
    "file_path": "path/file-b.txt",
    "file_name": "file-b.txt",
    "content": "x\\ny\\nz"
  }
]
```

The data source is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

## Configuration

The data source doesn't support configuration.

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data txt {
  # A glob pattern to select TXT files for reading
  #
  # For example:
  # glob = "path/to/files*.txt"
  #
  # Optional string. Default value:
  glob = null

  # A file path to a TXT file to read
  #
  # For example:
  # path = "data/disclaimer.txt"
  #
  # Optional string. Default value:
  path = null
}
```