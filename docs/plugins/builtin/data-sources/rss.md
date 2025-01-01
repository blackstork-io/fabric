---
title: "`rss` data source"
plugin:
  name: blackstork/builtin
  description: "Fetches RSS / Atom / JSON feed from a provided URL"
  tags: ["rss","http"]
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.2" "rss" "data source" >}}

## Description
Fetches RSS / Atom / JSON feed from a provided URL.

The full content of the items can be fetched and added to the feed. The data source supports basic authentication.

The data source is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

## Configuration

The data source doesn't support any configuration arguments.

## Usage

The data source supports the following execution arguments:

```hcl
data rss {
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


  # Required string.
  # For example:
  url = "https://www.elastic.co/security-labs/rss/feed.xml"

  # If the full content should be added when it's not present in the feed items.
  #
  # Optional bool.
  # Default value:
  fill_in_content = false

  # Maximum number of items to fill the content in per feed.
  #
  # Optional number.
  # Must be >= 0
  # For example:
  # fill_in_max_items = false
  # 
  # Default value:
  fill_in_max_items = 10

  # Return only items after a specified date time, in the format "%Y-%m-%dT%H:%M:%S%Z".
  #
  # Optional string.
  # For example:
  # only_items_after_time = "2024-12-23T00:00:00Z"
  # 
  # Default value:
  only_items_after_time = null
}
```