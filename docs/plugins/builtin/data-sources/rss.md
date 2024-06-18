---
title: "`rss` data source"
plugin:
  name: blackstork/builtin
  description: "Fetches RSS / Atom feed from a URL"
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
Fetches RSS / Atom feed from a URL.

The data source supports basic authentication.

The data source is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

## Configuration

The data source doesn't support any configuration arguments.

## Usage

The data source supports the following execution arguments:

```hcl
data rss {
  # Required string.
  # For example:
  url = "https://www.elastic.co/security-labs/rss/feed.xml"

  # Basic authentication credentials to be used in a HTTP request fetching RSS feed.
  #
  # Optional
  basic_auth {
    # Required string.
    # For example:
    username = "user@example.com"

    # Note: avoid storing credentials in the templates. Use environment variables instead.
    #
    # Required string.
    # For example:
    password = "passwd"
  }
}
```