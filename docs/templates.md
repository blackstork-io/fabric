---
title: Templates
description: Discover a diverse range of data sources within Fabric plugins. These integrations empower you to effortlessly load data from files, external services, APIs, and data storage systems. Simplify your data retrieval process and enhance your document generation workflow with Fabric's versatile data sources.
type: docs
weight: 60
---

# Templates

Fabric templates are the blueprints for configuring and generating documents. They offer a structured framework for defining data requirements, content structures, and rendering specifications within the [Fabric Configuration Language]({{< ref "language" >}}).

With Fabric templates, users can streamline their document generation process, ensuring consistency, accuracy, and scalability.

<div class="relative h-[35rem] w-auto not-prose">
  <img src="/images/fabric-template-example.png" alt="Fabric template code" title="Fabric template code" class="absolute left-0 top-0 z-0 w-[600px]"/>
  <img src="/images/fabric-template-example-rendered.png" alt="Rendered Fabric template" title="Rendered Fabric template" class="absolute right-0 bottom-0 z-10 w-[500px]"/>
</div>

## Templates repository

[Fabric Templates](https://github.com/blackstork-io/fabric-templates) repository contains use-case focused open-source templates for Fabric and is a great place to start:

- Cyber Threat Intelligence:
  - MITRE CTID CTI Blueprints ([source](https://mitre-engenuity.org/cybersecurity/center-for-threat-informed-defense/our-work/cti-blueprints/))
    - [Campaign Report Template](https://github.com/blackstork-io/fabric-templates/tree/main/cybersec/cti/mitre-ctid-campaign-report.fabric) ([example](https://github.com/blackstork-io/fabric-templates/tree/main/cybersec/cti/mitre-ctid-campaign-report-example.md))
    - [Executive Report Template](https://github.com/blackstork-io/fabric-templates/tree/main/cybersec/cti/mitre-ctid-executive-report.fabric) ([example](https://github.com/blackstork-io/fabric-templates/tree/main/cybersec/cti/mitre-ctid-executive-report-example.md))
    - [Intrusion Analysis Report Template](https://github.com/blackstork-io/fabric-templates/tree/main/cybersec/cti/mitre-ctid-intrusion-analysis-report.fabric) ([example](https://github.com/blackstork-io/fabric-templates/tree/main/cybersec/cti/mitre-ctid-intrusion-analysis-report-example.md))
    - [Threat Actor Profile Report Template](https://github.com/blackstork-io/fabric-templates/tree/main/cybersec/cti/mitre-ctid-threat-actor-profile-report.fabric) ([example](https://github.com/blackstork-io/fabric-templates/tree/main/cybersec/cti/mitre-ctid-threat-actor-profile-report-example.md))
- SecOps:
  - [Weekly Activity Overview Template](https://github.com/blackstork-io/fabric-templates/tree/main/cybersec/secops/weekly-activity-overview-elastic-security.fabric) ([example](https://github.com/blackstork-io/fabric-templates/tree/main/cybersec/secops/weekly-activity-overview-elastic-security-example.md))
