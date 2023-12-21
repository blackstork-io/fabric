package main

import (
	"github.com/hashicorp/hcl/v2"
)

const (
	ContentBlockName  = "content"
	DataBlockName     = "data"
	DocumentBlockName = "document"
)

type Templates struct {
	ContentBlocks []ContentBlock `hcl:"content,block"`
	DataBlocks    []DataBlock    `hcl:"data,block"`
	Documents     []Document     `hcl:"document,block"`
}

type MetaBlock struct {
	Name        *string   `hcl:"name,optional"`
	Author      *string   `hcl:"author,optional"`
	Description *string   `hcl:"description,optional"`
	Tags        []*string `hcl:"tags,optional"`
	UpdatedAt   *string   `hcl:"updated_at,optional"`

	RequiredFields []*string `hcl:"required_fields,optional"`
}

type DataBlock struct {
	Type    string `hcl:"type,label"`
	Name    string `hcl:"type,label"`
	Attrs   hcl.Attributes
	Meta    *MetaBlock `hcl:"meta,block"`
	Decoded bool
	Extra   hcl.Body `hcl:",remain"`
}

type DataBlockExtra struct {
	Ref   hcl.Expression `hcl:"ref,optional"`
	Extra hcl.Body       `hcl:",remain"`
}

type ContentBlock struct {
	Type    string `hcl:"type,label"`
	Name    string `hcl:"name,label"`
	Attrs   hcl.Attributes
	Meta    *MetaBlock `hcl:"meta,block"`
	Decoded bool

	Query *string `hcl:"query,optional"`
	Title *string `hcl:"title,optional"`

	Unparsed            hcl.Body `hcl:",remain"`
	NestedContentBlocks []ContentBlock

	localDict map[string]any
}

type ContentBlockExtra struct {
	Ref           hcl.Expression `hcl:"ref,optional"`
	ContentBlocks []ContentBlock `hcl:"content,block"`
	Unparsed      hcl.Body       `hcl:",remain"`
}

type Document struct {
	Name string `hcl:"name,label"`

	Meta  *MetaBlock `hcl:"meta,block"`
	Title *string    `hcl:"title,optional"`

	DataBlocks    []DataBlock    `hcl:"data,block"`
	ContentBlocks []ContentBlock `hcl:"content,block"`
}

// Block interfaces

type Block interface {
	// Data, common to all block kinds
	GetType() *string
	GetName() string
	GetAttrs() *hcl.Attributes
	GetMeta() **MetaBlock
	GetDecoded() *bool
	GetUnparsed() hcl.Body

	GetBlockKind() string
	// Get the structure for parsing this block's Extra fields
	NewBlockExtra() BlockExtra
	DecodeNestedBlocks(decoder *Decoder, extraInfo BlockExtra) hcl.Diagnostics
	UpdateFromRef(refTgt any, ref hcl.Expression) hcl.Diagnostics
}

type BlockExtra interface {
	GetRef() hcl.Expression
	GetUnparsed() hcl.Body
}
