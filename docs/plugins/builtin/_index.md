---
title: Built-in
weight: 10
plugin:
  name: blackstork/builtin
  description: ""
  tags: []
  version: "v0.0.0-dev"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
type: docs
---

{{< plugin-header "blackstork/builtin" "builtin" "v0.0.0-dev" >}}

`fabric` binary includes a set of built-in data sources and content providers, available out-of-the-box.

## Data sources

- [`csv`]({{< relref "./data-sources/csv" >}})

- [`inline`]({{< relref "./data-sources/inline" >}})

- [`json`]({{< relref "./data-sources/json" >}})

- [`txt`]({{< relref "./data-sources/txt" >}})

## Content providers

- [`frontmatter`]({{< relref "./content-providers/frontmatter" >}})
- [`image`]({{< relref "./content-providers/image" >}})
- [`list`]({{< relref "./content-providers/list" >}})
- [`table`]({{< relref "./content-providers/table" >}})
- [`text`]({{< relref "./content-providers/text" >}})
- [`toc`]({{< relref "./content-providers/toc" >}})