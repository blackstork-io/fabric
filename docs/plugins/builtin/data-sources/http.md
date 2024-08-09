---
title: "`http` data source"
plugin:
  name: blackstork/builtin
  description: "Loads data from a URL"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/builtin/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/builtin" "builtin" "v0.4.2" "http" "data source" >}}

## Description
Loads data from a URL.

At the moment, the data source accepts only responses with UTF-8 charset and parses only responses
with MIME types `text/csv` or `application/json`.

If MIME type of the response is `text/csv` or `application/json`, the response
content will be parsed and returned as a JSON structure (similar to the behaviour of CSV and JSON data
sources). Otherwise, the response content will be returned as text

The data source is built-in, which means it's a part of `fabric` binary. It's available out-of-the-box, no installation required.

## Configuration

The data source doesn't support any configuration arguments.

## Usage

The data source supports the following execution arguments:

```hcl
data http {
  # Basic authentication credentials to be used for HTTP request.
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


  # URL to fetch data from. Supported schemas are `http` and `https`
  #
  # Required string.
  # Must be non-empty
  # For example:
  url = "https://example.localhost/file.json"

  # HTTP method for the request. Allowed methods are `GET`, `POST` and `HEAD`
  #
  # Optional string.
  # Must be one of: "GET", "POST", "HEAD"
  # Default value:
  method = "GET"

  # If set to `true`, disabled verification of the server's certificate.
  #
  # Optional bool.
  # Default value:
  insecure = false

  # The duration of a timeout for a request. Accepts numbers, with optional fractions and a unit suffix. For example, valid values would be: 1.5s, 30s, 2m, 2m30s, or 1h
  #
  # Optional string.
  # Default value:
  timeout = "30s"

  # The headers to be set in a request
  #
  # Optional map of string.
  # Default value:
  headers = null

  # Request body
  #
  # Optional string.
  # Default value:
  body = null
}
```