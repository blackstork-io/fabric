---
title: Plugins
type: docs
weight: 5
---

# Plugins

[TBD]

Fabric relies on plugins called providers to interact with cloud providers, SaaS providers, and other APIs.

Fabric configurations must declare which providers they require so that Fabric can install and use them. Additionally, some providers require configuration (like endpoint URLs or cloud regions) before they can be used.

What Providers Do
Each provider adds a set of resource types and/or data sources that Fabric can manage.

Every resource type is implemented by a provider; without providers, Fabric can't manage any kind of infrastructure.

Most providers configure a specific infrastructure platform (either cloud or self-hosted). Providers can also offer local utilities for tasks like generating random numbers for unique resource names.


Where Providers Come From
Providers are distributed separately from Fabric itself, and each provider has its own release cadence and version numbers.

The Fabric Registry is the main directory of publicly available Fabric providers, and hosts providers for most major infrastructure platforms.
