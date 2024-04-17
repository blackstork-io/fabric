---
title: json
plugin:
  name: blackstork/builtin
  description: "Imports and parses the files matching \"glob\""
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
Imports and parses the files matching "glob".
Results are presented using the following structure:
```json
  [
    {
      "filename": "<name of the file matched by glob>",
      "contents": {
        "contents of the file": "parsed as json"
      },
    },
    {
      "filename": "<next file>",
      "contents": {
        "next": "contents"
      },
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
  # A pattern that selects the json files to be read
  #
  # Required string. For example:
  glob = "reports/*_data.json"
}
```