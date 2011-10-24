package quadtree

import (
	"strconv"
	"fmt"
)

// Private interface for quadtrees.
// Implemented by both node and leaf.
type treeInt interface {
	View() *View
	//
	insert(x, y float64, elems []interface{}, p *treeInt)
	//
	survey(view []*View, fun func(x, y float64, e interface{}))
	//
	delete(view *View, pred func(x, y float64, e interface{}) bool, p *treeInt) (delCount int)
	//
	isEmptyLeaf() bool
	//
	Width() float64
	Height() float64
	String() string
}

// Returns a new empty QuadTree whose View extends from
// leftX to rightX across the x axis and
// topY down to bottomY along the y axis
// leftX < rightX
// topY < bottomY
func NewQuadTree(leftX, rightX, topY, bottomY float64) T {
	var newView = NewViewP(leftX, rightX, topY, bottomY)
	return newRoot(newView)
}

// A point with a slice of stored elements
type vpoint struct {
	x, y  float64
	elems []interface{}
}

// Indicates whether a vpoint is zeroed, i.e. uninitialised
func (np *vpoint) zeroed() bool {
	return np.elems == nil && np.x == 0 && np.y == 0
}

// Resets this vpoint back to its uninitialised state
func (np *vpoint) zeroOut() {
	np.elems = nil
	np.x = 0.0
	np.y = 0.0
}

// Indicates whether a vpoint has the same x,y coords as those passed in
func (np *vpoint) sameLoc(x, y float64) bool {
	return np.x == x && np.y == y
}

func (np *vpoint) String() string {
	x := strconv.Ftoa64(np.x, 'f', 3)
	y := strconv.Ftoa64(np.y, 'f', 3)
	return " (" + fmt.Sprint(np.elems) + ", " + x + ", " + y + ")"
}

const LEAF_SIZE = 16

// A leaf struct implements the interface treeInt. Like a node (see below),
// a leaf contains a View defining the rectangular area in which each vpoint
// could legally be located. A leaf struct may contain up to four non-zeroed
// vpoints.
// If any vpoint is initialised, return false to zeroed(), then all vpoints 
// of lesser index are also initialised i.e. if ps[3] is initialised then so
// are ps[2] and ps[1], while ps[4] has no such constraints.
// The vpoints are not ordered in any way with respect to their positions at
// the leaf level.
type leaf struct {
	view View
	ps   [LEAF_SIZE]vpoint
}

// Returns a new leaf struct containing view and with all vpoints zeroed.
func newLeaf(view *View) *leaf {
	return &leaf{view: *view}
}

// Returns a pointer to the View of this leaf
func (l *leaf) View() *View {
	return &l.view
}

// Inserts each of the elements in elems into this leaf. There are three
// cases.
// If:			We find an existing vpoint at the exact location (x,y) 
//					- Append elems to this vpoint
// Else-If:	We find a zeroed vpoint available 
//					- Append elems to this vpoint
// Else:		This leaf has overflowed 
//					- Replace this leaf with an intermediate node and re-allocate 
//					all of the elements in this leaf as well as those in elems into
//					the new node
func (l *leaf) insert(x, y float64, elems []interface{}, inPtr *treeInt) {
	for i := range l.ps {
		if l.ps[i].zeroed() {
			l.ps[i].x = x
			l.ps[i].y = y
			l.ps[i].elems = append(l.ps[i].elems, elems...)
			return
		}
		if l.ps[i].sameLoc(x, y) {
			l.ps[i].elems = append(l.ps[i].elems, elems...)
			return
		}
	}
	// This leaf is full we need to create an intermediary node to divide it up
	newIntNode(x, y, elems, inPtr, l)
}

func newIntNode(x, y float64, elems []interface{}, inPtr *treeInt, l *leaf) {
	var intNode treeInt
	intNode = newNode(&l.view)
	*inPtr = intNode // Redirect the old leaf's reference to this intermediate node
	for _, p := range l.ps {
		intNode.insert(p.x, p.y, p.elems, &intNode)
	}
	intNode.insert(x, y, elems, &intNode)
}

// Applies fun to each of the elements contained in this leaf
// which appear within view.
func (l *leaf) survey(vs []*View, fun func(x, y float64, e interface{})) {
	for i := range l.ps {
		p := &l.ps[i]
		if !p.zeroed() && contains(vs, p.x, p.y) {
			for i := range p.elems {
				fun(p.x, p.y, p.elems[i])
			}
		}
	}
}

// Deletes each element, e, in this leaf which satisfies two conditions
// 	1: e lies within view
//	2: pred(e) returns true
// It is expected that pred will have side-effects allowing for the processing of deleted elements.
func (l *leaf) delete(view *View, pred func(x, y float64, e interface{}) bool, _ *treeInt) (delCount int) {
	for i := range l.ps {
		point := &l.ps[i]
		if !point.zeroed() && view.contains(point.x, point.y) {
			delCount += delete(point, pred)
			if len(point.elems) == 0 {
				point.zeroOut()
			}
		}
	}
	restoreOrder(&l.ps)
	return
}

// Deletes each element, e, from elems where pred(e) returns true.
func delete(p *vpoint, pred func(x, y float64, e interface{}) bool) (delCount int) {
	for i := len(p.elems) - 1; i >= 0; i-- {
		if pred(p.x, p.y, p.elems[i]) {
			// Fast delete from slice
			last := len(p.elems) - 1
			p.elems[i] = p.elems[last]
			p.elems = p.elems[:last]
			delCount++
		}
	}
	return
}

// Restores the leaf invariant that "if any vpoint is initialised, then all vpoints 
// of lesser index are also initialised" by rearranging the elements of ps.
func restoreOrder(ps *[LEAF_SIZE]vpoint) {
	for i := range ps {
		if ps[i].zeroed() {
			for j := i + 1; j < len(ps); j++ {
				if !ps[j].zeroed() {
					ps[i] = ps[j]
					ps[j].zeroOut()
					break
				}
			}
		}
	}
}

// Indicates whether or not this leaf contains any elements
func (l *leaf) isEmptyLeaf() bool {
	return l.ps[0].zeroed()
}

func (l *leaf) Width() float64 {
	return l.view.width()
}

func (l *leaf) Height() float64 {
	return l.view.height()
}

func (l *leaf) String() string {
	var str = l.view.String()
	for _, p := range l.ps {
		str += p.String()
	}
	return str
}

// A node struct implements the treeInt interface.
// A node is the intermediate, non-leaf, storage structure for a 
// quadtree.
// It contains a View, indicating the rectangular area this node covers.
// Each subtree will have a view containing one of four quarters of
// this node's view. Every subtree is guaranteed to be non-nil and
// may be either a node or a leaf struct.
type node struct {
	view     View
	children [4]treeInt
}

// Private method for creating new nodes. Returns a node with
// four leaves within the View provided.
func newNode(view *View) *node {
	v1, v2, v3, v4 := view.quarters()
	n := &node{*view, [4]treeInt{newLeaf(v1), newLeaf(v2), newLeaf(v3), newLeaf(v4)}}
	return n
}

// Returns the View for this node
func (n *node) View() *View {
	return &n.view
}

func (n *node) insert(x, y float64, elems []interface{}, _ *treeInt) {
	for i := range n.children {
		if n.children[i].View().contains(x, y) {
			n.children[i].insert(x, y, elems, &n.children[i])
			return
		}
	}
}

func (n *node) survey(vs []*View, fun func(x, y float64, e interface{})) {
	for i := range n.children {
		if overlaps(vs, n.children[i].View()) {
			n.children[i].survey(vs, fun)
		}
	}
}

func (n *node) delete(view *View, pred func(x, y float64, e interface{}) bool, p *treeInt) (delCount int) {
	allEmpty := true
	for i := range n.children {
		if n.children[i].View().overlaps(view) {
			var nInt treeInt
			nInt = treeInt(n)
			delCount += n.children[i].delete(view, pred, &nInt)
		}
		allEmpty = allEmpty && n.children[i].isEmptyLeaf()
	}
	if allEmpty && p != nil {
		var l treeInt
		l = newLeaf(n.View())
		*p = l
	}
	return
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

// External Node - each tree
type root struct {
	size int
	n    node
}

// Returns a new root ready for use as an empty quadtree
func newRoot(v *View) *root {
	r := new(root)
	node := newNode(v)
	r.n = *node
	return r
}

// Inserts the value nval into this node
func (r *root) Insert(x, y float64, nval interface{}) {
	elems := make([]interface{}, 1, 1)
	elems[0] = nval
	r.n.insert(x, y, elems, nil)
	r.size++
}

// Deletes each element, e, under this node which satisfies two conditions
// 	1: e lies within view
//	2: pred(e) returns true
// It is expected that pred will have side-effects allowing for the processing of deleted elements.
func (r *root) Delete(view *View, pred func(x, y float64, e interface{}) bool) {
	r.size -= r.n.delete(view, pred, nil)
}

// Applies fun to every element occurring within view in this node
func (r *root) Survey(vs []*View, fun func(x, y float64, e interface{})) {
	r.n.survey(vs, fun)
}

// Returns the View for this node
func (r *root) View() *View {
	return &r.n.view
}

func (r *root) String() string {
	return r.n.String()
}
