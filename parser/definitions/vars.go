package definitions

import (
	"context"
	"slices"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
	"github.com/blackstork-io/fabric/plugin"
)

type Query interface {
	Eval(ctx context.Context, dataCtx plugin.MapData) (result plugin.Data, diags diagnostics.Diag)
	Range() *hcl.Range
}
type QueryResultPlaceholder int

var (
	QueryType                  = encapsulator.New[Query]("query")
	QueriesType                = encapsulator.New[Queries]("queries")
	QueryResultPlaceholderType = encapsulator.New[int]("query result placeholder")
)

// Key for the eval context map to store the pending queries array.
const QueryKey = "\x00queries"

type ParsedVars []*hclsyntax.Attribute

// Queries is a threadsafe collection of pending queries.
type Queries struct {
	mu      sync.Mutex
	queries []DeferredQuery
}

type DeferredQuery struct {
	Query Query
	// path to store the result of the query
	// elements may be int or string (array index or map key)
	ResultPath []any
}

// Append adds a query to the internal state and returns an id.
func (q *Queries) Append(query Query) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	idx := len(q.queries)
	q.queries = append(q.queries, DeferredQuery{Query: query})
	return idx
}

// Take returns all the queries and resets the internal state.
func (q *Queries) Take() []DeferredQuery {
	q.mu.Lock()
	defer q.mu.Unlock()
	queries := q.queries
	q.queries = nil
	return queries
}

// ResultDest sets the result destination for a query with `id`.
func (q *Queries) ResultDest(id int, path []any) {
	path = slices.Clone(path)
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queries[id].ResultPath = path
}

func NewQueries() *Queries {
	return &Queries{}
}
