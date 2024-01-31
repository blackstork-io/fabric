---
title: Refrence
type: docs
weight: 1
---

# Fabric Block Refrences

This page covers Data, Content, and Section **References**—essential tools for building a robust, reusable and adaptive configuration.

## Data Block References

In Fabric, the `ref` attribute is a helpfupl tool for cross-referencing data blocks, a crucial aspect of maintaining a coherent configuration. 
 
To reference this block in another document, use the [`Data`](#./blocks/data.md) module and specify the name:

```
data elasticsearch "foo" {
}

document "overview" {

  data ref {
    base = data.elasticsearch.foo
    ...
  }

  data ref "foo2" {
    base = data.elasticsearch.foo
    ...
  }

}
```

It's worth noting that any additional specifications in the referer block replace the original to add adaptability, so to keep the configuration structure clear, avoid nesting blocks within referrers.

## Content Block References

Moving onto a content example, the `openai` content block labeled "foo" can be seamlessly integrated elsewhere using the [`content`](#content) module:

```
content openai "foo" {
  ...
}

document "overview" {

  content ref {
    base = content.openai.foo
    ...
  }

}
```

The `base` attribute acts as the central hub, linking content blocks. Any enhancements made in the referer block supersede the original specifications, similar to data referencing.

## Section Block References

The section blocks let you group parts of your request. Particularly using [`section ref`](#section) can help with structure:

```
section "foo" {
  ...
}

document "overview" {

  section ref {
    base = section.foo
    ...
  }

}
```

Adding a title in the referer section block, concise and structured, takes precedence over the original. To maintain clarity, avoid nesting within section referrers.
