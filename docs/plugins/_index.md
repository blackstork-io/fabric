---
title: Plugins
description: Discover the power of Fabric plugins, which implement various data sources and content providers to enhance your templating experience. Data sources enable loading data from files, external services, APIs, or data storage, while content providers render document content locally or via external APIs, supporting text, tables, graphs, code, and more.
type: docs
weight: 50
---

# Plugins

## Data sources, content providers and publishers

Fabric plugins implement various [data sources]({{< ref "data-sources.md" >}}), [content providers]({{< ref "content-providers.md" >}}) and [publishers]({{< ref "publishers.md" >}}):

- **data sources** are integrations responsible for loading data from local or external sources:
  files, external services and API, databases or cloud storage solutions.
- **content providers** render the content locally or with external API (like LLM). The providers
  produce various types of content: text, tables, graphs, code, etc.
- **publishers** are outgoing integrations that deliver rendered document to local or external
  destinations, for storage or dissemination.

## Dependencies

The global configuration must include all required plugins (see [Global configuration]({{< ref "../language/configs.md/#global-configuration" >}})). A plugin name consists of a namespace (usually a vendor's name) and a short name. For example, [`blackstork/elastic`]({{< ref "./elastic/" >}}) plugin (built by BlackStork) implements [Elasticsearch data source]({{< ref "./elastic/data-sources/elasticsearch" >}}).

## Installation

Plugin releases are independent from Fabric releases and have their own release cycle and version.
You can find a list of plugins released by BlackStork at the [Releases page](https://github.com/blackstork-io/fabric/releases)
in Fabric GitHub repository.

To automatically download and install required plugins, use `fabric install` command. See [Installing plugins]({{< ref "../install.md#installing-plugins" >}}) for more details.

## Available plugins

{{< plugins >}}
