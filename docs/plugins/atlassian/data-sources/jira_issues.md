---
title: "`jira_issues` data source"
plugin:
  name: blackstork/atlassian
  description: "Retrieve issues from Jira"
  tags: []
  version: "v0.4.2"
  source_github: "https://github.com/blackstork-io/fabric/tree/main/internal/atlassian/"
resource:
  type: data-source
type: docs
---

{{< breadcrumbs 2 >}}

{{< plugin-resource-header "blackstork/atlassian" "atlassian" "v0.4.2" "jira_issues" "data source" >}}

## Description
Retrieve issues from Jira.

## Installation

To use `jira_issues` data source, you must install the plugin `blackstork/atlassian`.

To install the plugin, add the full plugin name to the `plugin_versions` map in the Fabric global configuration block (see [Global configuration]({{< ref "configs.md#global-configuration" >}}) for more details), as shown below:

```hcl
fabric {
  plugin_versions = {
    "blackstork/atlassian" = ">= v0.4.2"
  }
}
```

Note the version constraint set for the plugin.

## Configuration

The data source supports the following configuration arguments:

```hcl
config data jira_issues {
  # Account Domain.
  #
  # Required string.
  # Must be non-empty
  # For example:
  domain = "some string"

  # Account Email.
  #
  # Required string.
  # Must be non-empty
  # For example:
  account_email = "some string"

  # API Token.
  #
  # Required string.
  # Must be non-empty
  # For example:
  api_token = "some string"
}
```

## Usage

The data source supports the following execution arguments:

```hcl
data jira_issues {
  # Use expand to include additional information about issues in the response.
  #
  # Optional string.
  # Must be one of: "renderedFields", "names", "schema", "changelog"
  # For example:
  # expand = "names"
  # 
  # Default value:
  expand = null

  # A list of fields to return for each issue.
  #
  # Optional list of string.
  # For example:
  # fields = ["*all"]
  # 
  # Default value:
  fields = null

  # A JQL expression. For performance reasons, this field requires a bounded query. A bounded query is a query with a search restriction.
  #
  # Optional string.
  # For example:
  # jql = "order by key desc"
  # 
  # Default value:
  jql = null

  # A list of up to 5 issue properties to include in the results.
  #
  # Optional list of string.
  # Must have a length of at most 5
  # Default value:
  properties = []

  # Size limit to retrieve.
  #
  # Optional number.
  # Must be >= 0
  # Default value:
  size = 0
}
```