package eval

import (
	"context"
	"fmt"
	"sync"

	"iter"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

type tuple[T1, T2 any] struct {
	V1 T1
	V2 T2
}

func fromVals[T1, T2 any](v1 T1, v2 T2) tuple[T1, T2] {
	return tuple[T1, T2]{v1, v2}
}

func (t tuple[T1, T2]) Vals() (T1, T2) {
	return t.V1, t.V2
}

func mappedInParallel[I any, C ~[]I, O1 any, O2 any](ctx context.Context, input C, f func(context.Context, I) (O1, O2)) iter.Seq2[int, tuple[O1, O2]] {
	var mu sync.Mutex
	remaining := len(input)
	done := remaining == 0
	ctx, cancel := context.WithCancel(ctx)
	if done {
		cancel()
	}
	return func(yield func(int, tuple[O1, O2]) bool) {
		for i, v := range input {
			go func() {
				t := fromVals(f(ctx, v))
				mu.Lock()
				defer mu.Unlock()
				if done {
					return
				}
				remaining--
				if !yield(i, t) || remaining == 0 {
					done = true
					cancel()
				}
			}()
		}
		<-ctx.Done()
		mu.Lock()
		defer mu.Unlock()
		done = true
	}
}

type customNodeRenderer struct {
	root              *nodes.Node
	publisherInfo     plugin.PublisherInfo
	supportedNodesSet map[string]struct{}
	allNodeRenderers  map[string]struct{}
	publisherName     string
	nodeRenderers     NodeRenderers
}

type customNodeRenderRequest struct {
	// subtree to be sent to the custom node renderer
	scope *nodes.Node
	// relative path to the custom node
	relativeNodePath nodes.Path
}

func newCustomNodeRenderer(root *nodes.Node, publisherName string, info plugin.PublisherInfo, nodeRenderers NodeRenderers) *customNodeRenderer {
	return &customNodeRenderer{
		root:              root,
		publisherInfo:     info,
		supportedNodesSet: utils.SliceToSet(info.SupportedCustomNodes),
		nodeRenderers:     nodeRenderers,
		allNodeRenderers:  nodeRenderers.AllNodeRenderers(),
		publisherName:     publisherName,
	}
}

func (r *customNodeRenderer) render(ctx context.Context, req *customNodeRenderRequest) (repl []*nodes.Node, diags diagnostics.Diag) {
	customNode := req.scope.TraversePath(req.relativeNodePath)
	content, ok := customNode.Content.(*nodes.Custom)
	if !ok {
		diags.Add(
			"Incorrect node type",
			fmt.Sprintf("Expected custom node, got %T", customNode.Content),
		)
		return
	}

	renderer, found := r.nodeRenderers.NodeRenderer(content.Data.GetTypeUrl())
	if !found {
		diags.Add(
			"Custom node renderer not found",
			fmt.Sprintf("Custom node renderer for type '%s' not found", content.Data.GetTypeUrl()),
		)
		return
	}
	var diag diagnostics.Diag
	repl, diag = renderer(ctx, &plugin.RenderNodeParams{
		Subtree:         req.scope,
		NodePath:        req.relativeNodePath,
		Publisher:       r.publisherName,
		PublisherInfo:   r.publisherInfo,
		CustomRenderers: r.allNodeRenderers,
	})
	diags.Extend(diag)
	return
}

// renderNodes renders all custom nodes in the tree
func (r *customNodeRenderer) renderNodes(ctx context.Context) (root *nodes.Node, diags diagnostics.Diag) {
	reqs := r.scan()
	for renderStep := 0; renderStep < 1000 && len(reqs) != 0; renderStep++ {
		// submit render requests in parallel
		for i, tuple := range mappedInParallel(ctx, reqs, r.render) {
			repl, diag := tuple.Vals()
			if diags.Extend(diag) {
				// if some error occurred - delete the custom node from the tree
				reqs[i].scope.TraversePath(reqs[i].relativeNodePath).RemoveFromTree()
				continue
			}
			// update the tree with the result
			if reqs[i].scope.Parent() == nil {
				reqs[i].scope.ReplaceWith(repl...)
				continue
			}

			// replacing the root node
			if len(repl) != 1 {
				diags.Add(
					"Invalid root node replacement",
					fmt.Sprintf("Root node replacement must be a single node, got %d nodes", len(repl)),
				)
				reqs[i].scope.TraversePath(reqs[i].relativeNodePath).RemoveFromTree()
				continue
			}
			newRoot := repl[0]
			if _, ok := newRoot.Content.(*nodes.FabricDocument); !ok {
				diags.Add(
					"Invalid root node replacement",
					"Root node replacement must be a document node",
				)
				reqs[i].scope.TraversePath(reqs[i].relativeNodePath).RemoveFromTree()
				continue
			}
			r.root = newRoot
		}
		reqs = r.scan()
	}
	if len(reqs) != 0 {
		diags.AddWarn(
			"Custom node rendering took too long",
			"Custom node rendering took too long, possibly due to circular dependencies",
		)
	}
	root = r.root
	return
}

// emptyNode is used as a placeholder
var emptyNode = new(nodes.Node)

// scan produces a list of independent custom node render requests
func (r *customNodeRenderer) scan() (reqs []*customNodeRenderRequest) {
	var contextNodes []*nodes.Node
	scopeIdxs := [nodes.ScopeCount]int{-1, -1, -1, -1}
	var stack []int

	pushScope := func(scope nodes.RendererScope) {
		stack = append(stack, scopeIdxs[scope])
		scopeIdxs[nodes.ScopeDocument] = len(contextNodes)
		contextNodes = append(contextNodes, emptyNode)
	}
	popScope := func(scope nodes.RendererScope, cursor *nodes.Cursor) {
		scopeIdxs[scope], stack = utils.PopSlice(stack)
		var n *nodes.Node
		n, contextNodes = utils.PopSlice(contextNodes)
		if n != nil && n != emptyNode {
			// have a valid scoped node to render
			reqs = append(reqs, &customNodeRenderRequest{
				scope:            cursor.Node(),
				relativeNodePath: n.GetRelativePath(cursor.Node()),
			})
		}
	}

	for cursor := range r.root.Walk() {
		switch c := cursor.Content().(type) {
		case *nodes.FabricDocument:
			if cursor.IsEntering() {
				pushScope(nodes.ScopeDocument)
				pushScope(nodes.ScopeSection)
			} else {
				popScope(nodes.ScopeSection, cursor)
				popScope(nodes.ScopeDocument, cursor)
			}
		case *nodes.FabricSection:
			if cursor.IsEntering() {
				pushScope(nodes.ScopeSection)
			} else {
				popScope(nodes.ScopeSection, cursor)
			}
		case *nodes.FabricContent:
			if cursor.IsEntering() {
				pushScope(nodes.ScopeContent)
			} else {
				popScope(nodes.ScopeContent, cursor)
			}
		case *nodes.Custom:
			if cursor.IsEntering() {
				continue
			}
			// TODO: check that the type is registered
			if utils.Contains(r.supportedNodesSet, c.Data.GetTypeUrl()) {
				// if the node is natively supported by the publisher - leave it as is
				continue
			}

			nodeIdx := scopeIdxs[c.Scope]
			if nodeIdx == -1 {
				// if scope is missing - treat as if the scope is Node
				reqs = append(reqs, &customNodeRenderRequest{
					scope:            cursor.Node(),
					relativeNodePath: nodes.Path{},
				})
				clear(contextNodes)
			} else {
				if contextNodes[nodeIdx] == emptyNode {
					contextNodes[nodeIdx] = cursor.Node()
					clear(contextNodes[:nodeIdx])
				}
				break
			}
		}
	}
	return reqs
}
