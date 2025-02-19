---
title: "`sleep` data source"
plugin:
  name: blackstork/builtin
  description: "Sleeps for the specified duration. Useful for testing and debugging"
  tags: ["debug"]
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.2" "sleep" "data source" >}}

## Description

Sleeps for the specified duration. Useful for testing and debugging.


The data source is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

## Configuration

The data source doesn't support any configuration arguments.

## Usage

The data source supports the following execution arguments:

```hcl
data sleep {
  # Duration to sleep
  #
  # Optional string.
  # Must be non-empty
  # Default value:
  duration = "1s"
}
```