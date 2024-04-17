---
title: rss
plugin:
  name: blackstork/builtin
  description: "Fetches an rss or atom feed"
  tags: ["rss","http"]
  version: "v0.4.1"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.1" "rss" "data source" >}}

## Description
Fetches an rss or atom feed

The data source is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

## Configuration

The data source supports the following configuration parameters:

```hcl
config data rss {
  # Authentication parameters used while accessing the rss source.
  #
  # Optional
  basic_auth {
    # Required string. For example:
    username = "user@example.com"

    # Note: you can use function like "from_env()" to avoid storing credentials in plaintext
    #
    # Required string. For example:
    password = "passwd"
  }
}
```

## Usage

The data source supports the following parameters in the data blocks:

```hcl
data rss {
  # Required string. For example:
  url = "https://www.elastic.co/security-labs/rss/feed.xml"
}
```