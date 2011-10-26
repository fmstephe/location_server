package quadtree

import (
	"strconv"
	"fmt"
	"os"
)

var noLeaves = os.NewError("No leaves available")
var noNodes = os.NewError("No nodes available")

// Private interface for quadtrees.
// Implemented by both node and leaf.
type treeInt interface {
	View() *View
	//
	insert(x, y float64, elems []interface{}, p *treeInt, r *root) (err os.Error)
	//
	survey(view []*View, fun func(x, y float64, e interface{}))
	//
	delete(view *View, pred func(x, y float64, e interface{}) bool, p *treeInt, r *root)
	//
	isEmptyLeaf() bool
	//
	setView(view *View)
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
func NewQuadTree(leftX, rightX, topY, bottomY float64, minLeafMax int64) T {
	var newView = NewViewP(leftX, rightX, topY, bottomY)
	return newRoot(newView, minLeafMax)
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

const LEAF_SIZE = 4

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
	nextFree *leaf
	view     View
	ps       [LEAF_SIZE]vpoint
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
func (l *leaf) insert(x, y float64, elems []interface{}, inPtr *treeInt, r *root) (err os.Error) {
	for i := range l.ps {
		if l.ps[i].zeroed() {
			l.ps[i].x = x
			l.ps[i].y = y
			l.ps[i].elems = elems // append(l.ps[i].elems, elems...)
			println("\tInserting", l.view.String(), x, y, l, elems)
			return
		}
		if l.ps[i].sameLoc(x, y) {
			l.ps[i].elems = append(l.ps[i].elems, elems...)
			return
		}
	}
	// This leaf is full we need to create an intermediary node to divide it up
	return newIntNode(x, y, elems, inPtr, l, r)
}

func newIntNode(x, y float64, elems []interface{}, inPtr *treeInt, l *leaf, r *root) (err os.Error) {
	println("New Internal")
	var intNode treeInt
	intNode, err = r.newNode(l.View())
	if err != nil {
		return
	}
	for _, p := range l.ps {
		intNode.insert(p.x, p.y, p.elems, &intNode, r)
	}
	println("Adding Extra", x, y, elems)
	println(intNode.(*node), (*inPtr).(*node))
	intNode.insert(x, y, elems, &intNode, r)
	r.recycle(*inPtr)
	*inPtr = intNode // Redirect the old leaf's reference to this intermediate node
	return
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
func (l *leaf) delete(view *View, pred func(x, y float64, e interface{}) bool, _ *treeInt, _ *root) {
	for i := range l.ps {
		point := &l.ps[i]
		if !point.zeroed() && view.contains(point.x, point.y) {
			delete(point, pred)
			if len(point.elems) == 0 {
				point.zeroOut()
			}
		}
	}
	restoreOrder(&l.ps)
	return
}

// Deletes each element, e, from elems where pred(e) returns true.
func delete(p *vpoint, pred func(x, y float64, e interface{}) bool) {
	for i := len(p.elems) - 1; i >= 0; i-- {
		if pred(p.x, p.y, p.elems[i]) {
			// Fast delete from slice
			last := len(p.elems) - 1
			p.elems[i] = p.elems[last]
			p.elems = p.elems[:last]
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

func (l *leaf) setView(view *View) {
	l.view = *view
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
	nextFree *node
	view     View
	children [4]treeInt
}

// Returns the View for this node
func (n *node) View() *View {
	return &n.view
}

func (n *node) insert(x, y float64, elems []interface{}, _ *treeInt, r *root) (err os.Error) {
	var thisNode treeInt
	thisNode = n
	for i := range n.children {
		if n.children[i].View().contains(x, y) {
			child := n.children[i]
			if err := child.insert(x, y, elems, &thisNode, r); err != nil {
				return err
			}
			return
		}
	}
	return
}

func (n *node) survey(vs []*View, fun func(x, y float64, e interface{})) {
	for i := range n.children {
		child := n.children[i]
		if overlaps(vs, child.View()) {
			n.children[i].survey(vs, fun)
		}
	}
}

func (n *node) delete(view *View, pred func(x, y float64, e interface{}) bool, p *treeInt, r *root) {
	allEmpty := true
	for i := range n.children {
		if n.children[i].View().overlaps(view) {
			var nInt treeInt
			nInt = treeInt(n)
			n.children[i].delete(view, pred, &nInt, r)
		}
		allEmpty = allEmpty && n.children[i].isEmptyLeaf()
	}
	if allEmpty && p != nil {
		var l treeInt
		l, _ = r.newLeaf(n.View()) // TODO Think hard about whether this could error out
		*p = l
	}
	return
}

func (n *node) setView(view *View) {
	n.view = *view
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
	freeNode *node
	freeLeaf *leaf
	leaves   []leaf
	nodes    []node
	rootNode *node
}

// Returns a new root ready for use as an empty quadtree
func newRoot(view *View, minLeafMax int64) *root {
	leafNum := 3 - ((minLeafMax - 1) % 3) + minLeafMax
	nodeNum := (leafNum - 1) / 3
	r := new(root)
	r.leaves = make([]leaf, leafNum, leafNum)
	for i := 0; i < len(r.leaves)-2; i++ {
		r.leaves[i].nextFree = &r.leaves[i+1]
	}
	r.nodes = make([]node, nodeNum, nodeNum)
	for i := 0; i < len(r.nodes)-2; i++ {
		r.nodes[i].nextFree = &r.nodes[i+1]
	}
	r.freeNode = &r.nodes[0]
	r.freeLeaf = &r.leaves[0]
	rootNode, _ := r.newNode(view) // TODO handle this error condition or prevent it from possibly happening
	r.rootNode = rootNode
	return r
}

func (r *root) recycle(ti treeInt) {
	switch t := ti.(type) {
	case *leaf:
		//println("Recycle Node")
		r.recycleLeaf(ti.(*leaf))
	case *node:
		//println("Recycle Leaf")
		r.recycleNode(ti.(*node))
	}
}

// Private method for creating new nodes. Returns a node with
// four leaves within the View provided.
func (r *root) newNode(view *View) (n *node, err os.Error) {
	n = r.freeNode
	if n == nil {
		err = noNodes
		return
	}
	r.freeNode = n.nextFree
	n.view = *view
	if err = r.newLeaves(view, &n.children); err != nil {
		r.freeNode = n
		n = nil
		return
	}
	println(n.view.String(),n.children[0].View().String())
	return
}

func (r *root) recycleNode(n *node) {
	n.nextFree = r.freeNode
	r.freeNode = n
	for i := range n.children { // You cannot recycle a node with nil children
		r.recycle(n.children[i])
	}
	n.children = *new([4]treeInt)
	n.view = *new(View)
}

func (r *root) newLeaf(view *View) (l *leaf, err os.Error) {
	l = r.freeLeaf
	if l == nil {
		err = noLeaves
		return
	}
	r.freeLeaf = l.nextFree
	l.view = *view
	return
}

func (r *root) newLeaves(pView *View, leaves *[4]treeInt) (err os.Error) {
	initFree := r.freeLeaf
	for i := range leaves {
		leaves[i] = r.freeLeaf
		r.freeLeaf = r.freeLeaf.nextFree
		if r.freeLeaf == nil {
			err = noLeaves
			break
		}
	}
	if err != nil {
		r.freeLeaf = initFree
		*leaves = *new([4]treeInt)
	} else {
		v0, v1, v2, v3 := pView.quarters()
		leaves[0].setView(v0)
		leaves[1].setView(v1)
		leaves[2].setView(v2)
		leaves[3].setView(v3)
	}
	return
}

func (r *root) recycleLeaf(l *leaf) {
	l.nextFree = r.freeLeaf
	r.freeLeaf = l
	l.view = *new(View)
	l.ps = *new([LEAF_SIZE]vpoint)
}
// Inserts the value nval into this node
func (r *root) Insert(x, y float64, nval interface{}) {
	println("Inserting", x, y, nval)
	elems := make([]interface{}, 1, 1)
	elems[0] = nval
	err := r.rootNode.insert(x, y, elems, nil, r)
	if err != nil {
		panic("Unable to insert element")
	}
}

// Deletes each element, e, under this node which satisfies two conditions
// 	1: e lies within view
//	2: pred(e) returns true
// It is expected that pred will have side-effects allowing for the processing of deleted elements.
func (r *root) Delete(view *View, pred func(x, y float64, e interface{}) bool) {
	r.rootNode.delete(view, pred, nil, r)
}

// Applies fun to every element occurring within view in this node
func (r *root) Survey(vs []*View, fun func(x, y float64, e interface{})) {
	r.rootNode.survey(vs, fun)
}

// Returns the View for this node
func (r *root) View() *View {
	return &r.rootNode.view
}

func (r *root) freeNodes() (cnt int) {
	freeNode := r.freeNode
	for {
		if freeNode != nil {
			freeNode = freeNode.nextFree
			cnt++
		} else {
			return
		}
	}
	panic("Unreachable")
}

func (r *root) freeLeaves() (cnt int) {
	freeLeaf := r.freeLeaf
	for {
		if freeLeaf != nil {
			freeLeaf = freeLeaf.nextFree
			cnt++
		} else {
			return
		}
	}
	panic("Unreachable")
}

func (r *root) String() string {
	return r.rootNode.String()
}
