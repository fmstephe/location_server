package quadtree

import (
	"strconv"
	"fmt"
)

var fullTreeErr = os.NewError("Tree Full")

type arraytree struct {
	view View
	free int64
	nodes []node
}

func (t *arraytree) View() *View {
	return &t.view
}

func (td *arraytree) getFree() (i int64) os.Error {
	i = td.free
	td.free = td.nodes[free].nextFree
	if td.free == -1 {
		td.free = i
		return fullTreeErr
	}
	return nil
}

// The three possible states for a node
const (
	empty = iota
	internal
	leaf
)

// A node struct implements the treeInt interface.
// A node is the intermediate, non-leaf, storage structure for a 
// quadtree.
// It contains a View, indicating the rectangular area this node covers.
// Each subtree will have a view containing one of four quarters of
// this node's view. Every subtree is guaranteed to be non-nil and
// may be either a node or a leaf struct.
type node struct {
	view    View
	indice	int64
	nextFree int64
	state	int
	// Internal node values
	children [4]int64
	// Leaf node values
	x, y float64
	elems []interface{}
}

// Private method for creating new nodes. Returns a node with
// four leaves within the View provided.
func newNode(view *View) *node {
	return &node{view: *view, state: empty, children: [4]int64{-1, -1, -1, -1}}
}

func (n *node) freeNode(td *arraytree) {
	n.nextFree = td.free
	td.free = n.indice
	n.makeEmpty()
}

func (n *node) makeEmpty() {
	n.state = free
	n.x, n.y = 0.0, 0.0
	n.elems = nil
	// We assume that n has no children
}

func (n *node) makeInternal(td *arraytree) {
	n.state = internal
	n.elems = nil
	n.x, n.y = 0.0, 0.0
	for i := 0; i < 4; i++ {
		n.children[i] = td.getFree()
	}
}

func (n *node) insert(x, y float64, elems []interface{}, parent int64, td *arraytree) os.Error {
	if !n.view.contains(x,y) {
		return
	}
	switch n.state {
	case empty:
		n.state = leaf
		n.x = x
		n.y = y
		n.elems = elems
	case internal:
		for i := range children {
			child := nodes[children[i]]
			if err := child.insert(x,y,elems,interface{},n.indice,nodes); err != nil {
				return err // Consistent state?
			}
		}
	case leaf :
		oX, oY := n.x, n.y
		oElems = n.elems
		if err := n.makeInternal(td); err != nil {
			return err // Does this leave the tree in a consistent state?
		}
		if err := n.insert(oX,oY,oElems,n.indice,td); err != nil {
			return err // consistent state?
		}
		if err := n.insert(x,y,elems,n.indice,td); err != nil {
			return err // Consistent state?
		}
	}
}

func (n *node) survey(vs []*View, fun func(x, y float64, e interface{}), td *arraytree) {
	if !overlaps(vs,n.view) {
		return
	}
	switch n.state {
	case empty:
		return
	case internal:
		for i := range children {
			child := nodes[children[i]]
			child.survey(vs, fun, td)
		}
	case leaf:
		for i := range n.elems {
			fun(n.x, n.y, n.elems[i])
		}
	}
}

func (n *node) delete(view *View, pred func(x, y float64, e interface{}) bool, parent int64, td *arraytree) (madeEmpty bool) {
	if !n.view.overlaps(view) {
		return false
	}
	switch n.state {
	case empty:
		return false
	case internal:
		allEmpty := true
		for i := range children {
			child := nodes[children[i]]
			allEmpty &= child.delete(view, pred, td)
		}
		if allEmpty {
			for i := range children {
				child := nodes[children[i]].freeNode(td)
				children[i] = -1
			}
		}
		return allEmpty
	case leaf:
		for i := range n.elems {
			if pred(n.x, n.y, n.elems[i]) {
				// Fast delete from slice
				last := len(n.elems) - 1
				n.elems[i] = n.elems[last]
				n.elems = n.elems[:last]
			}
		}
		return len(n.elems) == 0
	}
}

func (n *node) Width() float64 {
	return n.view.width()
}

func (n *node) Height() float64 {
	return n.view.height()
}

func (n *node) isEmptyLeaf() bool {
	return false
}

func (n *node) String() string {
	return "<" + n.view.String() + "-\n" + n.children[0].String() + ", \n" + n.children[1].String() + ", \n" + n.children[2].String() + ", \n" + n.children[3].String() + ">"
}
