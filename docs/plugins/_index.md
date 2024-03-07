---
title: Plugins
type: docs
weight: 60
---

# Plugins

Fabric relies on plugins for implementing integrations with data sources and content providers. The global configuration should specify all required plugins (see [Global configuration]({{< ref "../language/configs.md/#global-configuration" >}}) for the details). Additionally, some data sources and content providers themselves require configuration (for example, API keys, URLs, credentials, etc).

A plugin name consists of a namespace (a name of a plugin vendor) and a short name. For example, `blackstork/elastic` plugin built by [BlackStork](https://blackstork.io) implements Elasticsearch and Elastic Security Cases data sources.

## Where to get the plugins

Plugin releases are independent from Fabric releases. Plugins are distributed independently and have their own release cycle and version.
You can find a list of plugins released by BlackStork at the [Releases page](https://github.com/blackstork-io/fabric/releases) in Fabric GitHub.
