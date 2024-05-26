package plugin

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// DataPath is a type that represents a path in `Data`.
type DataPath struct {
	root string
	path []any
}

// NewDataPath creates a new DataPath. Specified path must not change after creation.
func NewDataPath(v []any) *DataPath {
	return &DataPath{path: v}
}

// Appends an index in the ListData, returns self for chaining.
func (p *DataPath) List(idx int) *DataPath {
	p.path = append(p.path, idx)
	return p
}

// Appends an index in the MapData, returns self for chaining.
func (p *DataPath) Map(key string) *DataPath {
	p.path = append(p.path, key)
	return p
}

// Clones the path.
func (p *DataPath) Clone() *DataPath {
	return &DataPath{path: slices.Clone(p.path), root: p.root}
}

// Clones the path.
func (p *DataPath) SetRootName(root string) *DataPath {
	p.root = root
	return p
}

// Truncates the path. Negative values remove n elements from the end.
func (p *DataPath) Truncate(newLen int) *DataPath {
	if newLen < 0 {
		newLen = len(p.path) + newLen
	}
	p.path = p.path[:newLen]
	return p
}

// String representation of the path.
func (p *DataPath) String() string {
	var sb strings.Builder
	sb.WriteString(p.root)

	for _, path := range p.path {
		switch p := path.(type) {
		case string:
			fmt.Fprintf(&sb, ".%q", p)
		case int:
			fmt.Fprintf(&sb, ".[%d]", p)
		default:
			panic("incorrect path element type")
		}
	}
	return sb.String()
}

// Gets element from the path.
func (p *DataPath) Get(data Data) (Data, diagnostics.Diag) {
	for i, part := range p.path {
		switch part := part.(type) {
		case string:
			m, ok := data.(MapData)
			if !ok {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Map expected",
					Detail:   fmt.Sprintf("Expected to find map at %s, got %T", p.Clone().Truncate(i+1), data),
				}}
			}
			if m == nil {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Map is nil",
					Detail:   fmt.Sprintf("Expected to index map at %s, but it's nil", p.Clone().Truncate(i+1).String()),
				}}
			}
			data, ok = m[part]
			if !ok {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Key not found",
					Detail:   fmt.Sprintf("Key %q not found in map at %s", part, p.Clone().Truncate(i+1).String()),
				}}
			}
		case int:
			l, ok := data.(ListData)
			if !ok {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "List expected",
					Detail:   fmt.Sprintf("Expected to find list at %s, got %T", p.Clone().Truncate(i+1), data),
				}}
			}
			if l == nil {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "List is nil",
					Detail:   fmt.Sprintf("Expected to index list at %s, but it's nil", p.Clone().Truncate(i+1).String()),
				}}
			}
			if part >= len(l) || part < 0 {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Index out of range",
					Detail:   fmt.Sprintf("Index %d is out of range of list at %s", part, p.Clone().Truncate(i+1).String()),
				}}
			}
			data = l[part]
		default:
			panic("incorrect path element type")
		}
	}
	return data, nil
}

func (p *DataPath) Set(data, val Data) (Data, diagnostics.Diag) {
	if len(p.path) == 0 {
		return val, nil
	}
	dataToChange, diags := p.Clone().Truncate(len(p.path) - 1).Get(data)
	if diags.HasErrors() {
		return data, diags
	}
	lastPart := p.path[len(p.path)-1]
	switch lastPart := lastPart.(type) {
	case string:
		m, ok := dataToChange.(MapData)
		if !ok {
			return data, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Map expected",
				Detail:   fmt.Sprintf("Expected to find map at %s, got %T", p, data),
			}}
		}
		if m == nil {
			return data, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Map is nil",
				Detail:   fmt.Sprintf("Expected to index map at %s, but it's nil", p),
			}}
		}
		m[lastPart] = val
	case int:
		l, ok := data.(ListData)
		if !ok {
			return data, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "List expected",
				Detail:   fmt.Sprintf("Expected to find list at %s, got %T", p.Clone().Truncate(-1), data),
			}}
		}
		if l == nil {
			return data, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "List is nil",
				Detail:   fmt.Sprintf("Expected to index list at %s, but it's nil", p.Clone().Truncate(-1)),
			}}
		}
		if lastPart >= len(l) || lastPart < 0 {
			return data, diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Index out of range",
				Detail:   fmt.Sprintf("Index %d is out of range of list at %s", lastPart, p.Clone().Truncate(-1)),
			}}
		}
		l[lastPart] = val
	}
	return data, nil
}
