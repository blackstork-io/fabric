---
title: Language
description: Fabric Configuration Language (FCL) drives Fabric's document generation capabilities. With FCL, users define data needs and template structures in .fabric files, streamlining content creation. Automate document production efficiently with Fabric and FCL.
type: docs
weight: 40
---

# Fabric configuration language

Fabric Configuration Language (FCL) serves as the core feature for Fabric, a powerful tool designed to streamline document generation. FCL enables users to express data requirements and template structures within configuration files using a lightweight syntax.

The document templates, defined in the configuration files, act as blueprints for consolidating data and creating Markdown documents. FCL empowers users to define, manage, and automate the document production process, delivering a sturdy and adaptable solution for content generation.

Fabric configuration files have the extension `.fabric` and contain configurations, data requirements, and content definitions. The Fabric configuration codebase may consist of many files and subdirectories.

## Core concepts

Building upon the [HashiCorp Configuration Language](https://github.com/hashicorp/hcl) (HCL), Fabric language shares similarities with the [Terraform Configuration Language](https://developer.hashicorp.com/terraform/language) and comprises two fundamental elements:

- **Blocks**: serve as containers defining objects, such as configurations, data requirements, or content structures. Blocks always include a block type and may have zero or more labels.
- **Arguments** assign values to names within blocks, facilitating the configuration process.

```hcl
# Named data block:

data elasticsearch "alerts" {
    index = ".alerts-security.alerts-*"
    query_string = "kibana.alert.severity:critical"
}

<BLOCK-TYPE> <PLUGIN> "<BLOCK-NAME>" {
    <ARGUMENT> = <VALUE>
}

# Anonymous configuration block for a data plugin:
config data elasticsearch {
    cloud_id = "my-elastic-cloud-id"
    api_key = "my-elastic-cloud-api-key"
}

<CONFIG-LABEL> <BLOCK-TYPE> <PLUGIN> {
    <ARGUMENT> = <VALUE>
}
```

See [Syntax](./syntax/) for more details on the FCL syntax.

## IDE support

Given that Fabric configuration language is built on HCL, IDE extensions designed for HCL syntax highlighting are applicable to Fabric files. It may be necessary to explicitly set the file type for `*.fabric` files to HCL.

For users of Microsoft Visual Studio Code, there is a dedicated [Fabric Extension for Visual Studio Code](https://github.com/blackstork-io/vscode-fabric) available, providing enhanced support for Fabric configurations within the IDE.

![A screenshot of Fabric Extension for Visual Studio Code](./vscode-fabric-screenshot.webp)
