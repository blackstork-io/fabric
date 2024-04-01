---
title: Plugins
description: Discover the power of Fabric plugins, which implement various data sources and content providers to enhance your templating experience. Data sources enable loading data from files, external services, APIs, or data storage, while content providers render document content locally or via external APIs, supporting text, tables, graphs, code, and more.
type: docs
weight: 50
hideChildren: true
---

# Plugins

Fabric plugins implement various [data sources]({{< ref "data-sources.md" >}}) and [content providers]({{< ref "content-providers.md" >}}):

- **data sources** are integrations responsible for loading data from the file, external service, API, or data storage.
- **content providers** render the document content locally or through the external API (like LLM). The content types include text, tables, graphs, code and other elements.

## Dependencies

The global configuration must include all required plugins (see [Global configuration]({{< ref "../language/configs.md/#global-configuration" >}})). A plugin name consists of a namespace (usually a vendor's name) and a short name. For example, [`blackstork/elastic`]({{< ref "./elastic/" >}}) plugin (built by BlackStork) implements [Elasticsearch data source]({{< ref "./elastic/data-sources/elasticsearch" >}}).

## Installation

Plugin releases are independent from Fabric releases and have their own release cycle and version. You can find a list of plugins released by BlackStork at the [Releases page](https://github.com/blackstork-io/fabric/releases) in Fabric GitHub repository.

Required plugins can be downloaded and installed automatically with `fabric install` command. See [Installing plugins]({{< ref "../install.md#installing-plugins" >}}) for more details.

## Available plugins

{{< plugins >}}
