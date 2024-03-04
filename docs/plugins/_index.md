---
title: Plugins
type: docs
weight: 60
---

# Plugins

Fabric relies on plugins to implement data sources and content providers. To utilise a plugin's data sources and content providers, it must be installed by Fabric. The global configuration should specify all required plugins (see [Global configuration]({{< ref "../language/configs.md/#global-configuration" >}}) for the details). Additionally, some data sources and content providers require configuration (for example, API keys, URLs, credentials, etc).

A plugin name consists of a namespace (a name of a plugin vendor) and a short name. For example, `blackstork/elasticsearch` plugin implements Elasticsearch client data source and is released by [BlackStork](https://blackstork.io).

## Where to get the plugins

Plugins are released and distributed independently from Fabric, with their own release cycle and version.
You can find a list of plugins released by BlackStork at the [Releases page](https://github.com/blackstork-io/fabric/releases) in Fabric GitHub.
