package quadtree

import (
	"os"
	"fmt"
)

var fullTreeErr = os.NewError("Tree Full")

type arraytree struct {
	view      View
	root      *arraynode
	firstFree *arraynode
	nodes     []arraynode
}

func NewArrayTree(leftX, rightX, topY, bottomY float64, maxSize int64) T {
	at := new(arraytree)
	at.view = NewView(leftX, rightX, topY, bottomY)
	at.nodes = make([]arraynode, maxSize, maxSize)
	for i := range at.nodes {
		at.nodes[i] = initArrayNode(int64(i))
	}
	for i := len(at.nodes) - 2; i >= 0; i-- {
		at.nodes[i].nextFree = &at.nodes[i+1]
	}
	at.root = &at.nodes[0]
	at.firstFree = at.root.nextFree
	at.root.state = emptyState
	at.root.view = at.view
	return at
}

func (at *arraytree) takeFree(freedNodes *[4]*arraynode) (err os.Error) {
	oFree := at.firstFree
	for i := range freedNodes {
		if at.firstFree == nil {
			err = fullTreeErr
			break
		}
		freedNodes[i] = at.firstFree
		at.firstFree = freedNodes[i].nextFree
		freedNodes[i].makeEmpty()
	}
	if err != nil {
		at.firstFree = oFree
		*freedNodes = *new([4]*arraynode)
	}
	return
}

func (at *arraytree) View() *View {
	return &at.view
}

func (at *arraytree) Insert(x, y float64, e interface{}) {
	at.root.insert(x, y, []interface{}{e}, at)
}

func (at *arraytree) Survey(vs []*View, fun func(x, y float64, e interface{})) {
	at.root.survey(vs, fun)
}

func (at *arraytree) Delete(v *View, pred func(x, y float64, e interface{}) bool) {
	at.root.delete(v, pred, at)
}

func (at *arraytree) String() string {
	return "<" + at.view.String() + "-\n" + at.root.treeString(at) + ">"
}

// The three possible states for a arraynode
const (
	freeState = iota
	emptyState
	internalState
	leafState
)

// A arraynode struct implements the treeInt interface.
// A arraynode is the intermediate, non-leaf, storage structure for a 
// quadtree.
// It contains a View, indicating the rectangular area this arraynode covers.
// Each subtree will have a view containing one of four quarters of
// this arraynode's view. Every subtree is guaranteed to be non-nil and
// may be either a arraynode or a leaf struct.
type arraynode struct {
	view     View
	nextFree *arraynode
	indice   int64
	state    int
	children [4]*arraynode
	// Leaf arraynode values
	x, y  float64
	elems []interface{}
}

// Private method for creating new arraynodes. Returns a arraynode with
// four leaves within the View provided.
func initArrayNode(indice int64) arraynode {
	return arraynode{indice: indice, state: freeState}
}

func (n *arraynode) makeFree(at *arraytree) {
	n.nextFree = at.firstFree
	at.firstFree = n
	n.state = freeState
	n.x, n.y = 0.0, 0.0
	n.elems = nil
}

func (n *arraynode) makeEmpty() {
	n.state = emptyState
	n.x, n.y = 0.0, 0.0
	n.elems = nil
	// We assume that n has no children
}

func (n *arraynode) makeInternal(at *arraytree) (err os.Error) {
	n.state = internalState
	n.elems = nil
	n.x, n.y = 0.0, 0.0
	err = at.takeFree(&n.children)
	if err == nil {
		v1, v2, v3, v4 := n.view.quarters()
		n.children[0].view = *v1
		n.children[1].view = *v2
		n.children[2].view = *v3
		n.children[3].view = *v4
	}
	return
}

func (n *arraynode) makeLeaf(x, y float64, elems []interface{}) {
	n.state = leafState
	n.x = x
	n.y = y
	n.elems = elems
}

func (n *arraynode) insert(x, y float64, elems []interface{}, at *arraytree) (err os.Error) {
	if !n.view.contains(x, y) {
		return
	}
	switch n.state {
	case emptyState:
		n.makeLeaf(x, y, elems)
	case internalState:
		for i := range n.children {
			child := n.children[i]
			if err := child.insert(x, y, elems, at); err != nil {
				return err // TODO Consistent state?
			}
		}
	case leafState:
		if n.x == x && n.y == y {
			n.elems = append(n.elems, elems...)
		} else {
			nX, nY := n.x, n.y
			nElems := n.elems
			if err := n.makeInternal(at); err != nil {
				return err // TODO Does this leave the tree in a consistent state?
			}
			if err := n.insert(nX, nY, nElems, at); err != nil {
				return err // TODO Consistent state?
			}
			if err := n.insert(x, y, elems, at); err != nil {
				return err // TODO Consistent state?
			}
		}
	}
	return
}

func (n *arraynode) survey(vs []*View, fun func(x, y float64, e interface{})) {
	if !overlaps(vs, &n.view) {
		return
	}
	switch n.state {
	case freeState:
		panic("Attempted to survey a free arraynode. Tree structure is broken.")
	case emptyState:
		return
	case internalState:
		for i := range n.children {
			child := n.children[i]
			child.survey(vs, fun)
		}
	case leafState:
		if contains(vs, n.x, n.y) {
			for i := range n.elems {
				fun(n.x, n.y, n.elems[i])
			}
		}
	}
}

func (n *arraynode) delete(view *View, pred func(x, y float64, e interface{}) bool, at *arraytree) (isEmpty bool) {
	if !n.view.overlaps(view) {
		return n.state == emptyState
	}
	switch n.state {
	case emptyState:
		return true
	case internalState:
		allEmpty := true
		for i := range n.children {
			emptyChild := n.children[i].delete(view, pred, at)
			allEmpty = allEmpty &&  emptyChild
		}
		if allEmpty {
			for i := range n.children {
				n.children[i].makeFree(at)
				n.children[i] = nil
			}
			n.makeEmpty()
		}
		return allEmpty
	case leafState:
		if view.contains(n.x, n.y) {
			for i := len(n.elems) - 1; i >= 0; i-- {
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
	return false
}

func (n *arraynode) treeString(at *arraytree) string {
	switch n.state {
	case emptyState:
		return "_"
	case internalState:
		out := "<"
		out += n.view.String() + ":"
		for i := range n.children {
			out += n.children[i].treeString(at)
			out += ","
		}
		out += ">"
		return out
	case leafState:
		return fmt.Sprintf("%v:%v", n.view, n.elems)
	}
	return "Illegal State"
}
