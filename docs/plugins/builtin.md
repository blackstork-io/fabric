---
title: Built-in
weight: 10
type: docs
---

# Built-in data sources and content providers

`fabric` binary includes a set of built-in data sources and content providers, available out-of-the-box.

## Data sources

### `csv`

#### Configuration

The data source supports the following configuration parameters:

```hcl
config data csv {
    delimiter = <string>  # optional
}
```

#### Usage

The data source supports the following parameters in the data blocks:

```hcl
data csv {
    path = <string>  # required
}
```

### `inline`

#### Configuration

The data source doesn't support configuration.

#### Usage

The data source doesn't define any parameters in the `data` block.

### `json`

#### Configuration

The data source doesn't support configuration.

#### Usage

The data source supports the following parameters in the data blocks:

```hcl
data json {
    glob = <string>  # required
}
```

### `txt`

#### Configuration

The data source doesn't support configuration.

#### Usage

The data source supports the following parameters in the data blocks:

```hcl
data txt {
    path = <string>  # required
}
```

## Content providers

### `frontmatter`

#### Configuration

The content provider doesn't support configuration.

#### Usage

The content source supports the following parameters in the content blocks:

```hcl
content frontmatter {
    content = <map of dynamic>  # optional
    format = <string>  # optional
}
```

### `image`

#### Configuration

The content provider doesn't support configuration.

#### Usage

The content source supports the following parameters in the content blocks:

```hcl
content image {
    alt = <string>  # optional
    src = <string>  # required
}
```

### `list`

#### Configuration

The content provider doesn't support configuration.

#### Usage

The content source supports the following parameters in the content blocks:

```hcl
content list {
    format = <string>  # optional
    item_template = <string>  # required
}
```

### `table`

#### Configuration

The content provider doesn't support configuration.

#### Usage

The content source supports the following parameters in the content blocks:

```hcl
content table {
    columns = <list of object>  # required
}
```

### `text`

#### Configuration

The content provider doesn't support configuration.

#### Usage

The content source supports the following parameters in the content blocks:

```hcl
content text {
    absolute_title_size = <number>  # optional
    code_language = <string>  # optional
    format_as = <string>  # optional
    text = <string>  # required
}
```

### `toc`

#### Configuration

The content provider doesn't support configuration.

#### Usage

The content source supports the following parameters in the content blocks:

```hcl
content toc {
    end_level = <number>  # optional
    ordered = <bool>  # optional
    scope = <string>  # optional
    start_level = <number>  # optional
}
```
