---
title: Language
type: docs
weight: 3
---

# Understanding Fabric Blocks
Blocks in Fabric stand as pivotal components, creating structure of your projects. Blocks prioritize readability, maintainability, and collaborative development, shaping  your processed documents.

Similar to Terraform, Blocks are one of the [two main components](../syntax/) that make up our language. Think of Blocks as blueprint of your document.

## Key Advantages of Fabric Blocks

Fabric Blocks offer substantial advantages, emphasizing:

- **Organization:** Efficiently structure configurations into manageable sections.
- **Collaboration:** Enable seamless teamwork by dissecting specific blocks.
- **Readability:** Foster clarity and comprehension through well-structured configurations.

## Types of Blocks

In Fabric, Blocks define the structure, content, and organization of your processed documents. Let's delve deeper into each type of block to understand their roles, syntax, and the problems they solve.

In Fabric there are three types of Blocks:
- Content
- Data
- Section

## Data Blocks

Data blocks defines a call to a data plugin, sourcing external data for content rendering. The order of data block definitions within the document is immaterial.

```
data <plugin-name> "<result-name>" {
  ...
}

document "foobar" {

  data <plugin-name> "<result-name>" {
    ...
  }

}
```

At the root level, ensure you provide both the plugin and result names; they serve as crucial identifiers within your codebase. If you're declaring data blocks outside the document, make sure to reference them inside your template.

Now, when defining data blocks within the document, you only need the plugin and result names, forming a unique pair that acts as identifiers. The parameters within the block serve as inputs for the data plugin, and the data generated is stored globally under data.<plugin-name>.<result-name>.

## Content Blocks

The content block defines a call to a content plugin that generates document segments. The order in which content blocks are defined is crucial, determining the sequence of generated content in the document.

```
content <plugin-name> "<block-name>" {
  ...
}

document "foobar" {

  content <plugin-name> "<block-name>" {
    ...
  }

  content <plugin-name> {
    ...
  }

}
```


At the root level, you'll need both the plugin name and, if desired, the block name, as they act as crucial identifiers within your codebase—essential for referencing the block. If you're defining content blocks outside the document, make sure to reference them within your template.

Now, when it comes to content blocks within the document, only the plugin's name is necessary; the block name is optional.

Within the block, you set the parameters, serving as inputs for the content plugin, alongside the plugin configuration and the local context map. These plugins are designed to return a Markdown text string.

And here's the thing – the order in which you define your content blocks? Yep, it's preserved.

## Section Blocks

Section blocks enable grouping content blocks for enhanced reusability and referencing. They also can contain other section blocks.

```
section "<section-name>" {
  
}

document "foobar" {
  section {
    ...
  }

  section "<section-name>" {
    section {
      ...
    }
    ...
  }
}
```

When specifying a section block at the root level, make sure to include a section name – a unique identifier crucial for referencing this block within your codebase. Section blocks declared outside the document need to be referenced inside the template.

However, if you're defining a section block within the document template, the section name becomes optional.

Embark on the journey of Fabric Blocks, where the dynamic interplay of structure and flexibility adds a sophisticated touch to your configuration landscape.