package quadtree

import (
	"fmt"
)

// Private interface for quadtree nodes. Implemented by both node and leaf.
type subtree interface {
	//
	insert(x, y float64, elems []interface{}, p *subtree, r *root)
	//
	survey(view []*View, fun func(x, y float64, e interface{}))
	//
	delete(view *View, pred func(x, y float64, e interface{}) bool, p *subtree, r *root)
	//
	isEmptyLeaf() bool
	//
	View() *View
	//
	setView(view *View)
	//
	String() string
}

// Returns a new empty QuadTree whose View extends from
// leftX to rightX across the x axis and
// topY down to bottomY along the y axis
// leftX < rightX
// topY < bottomY
func NewQuadTree(leftX, rightX, topY, bottomY float64, leafAllocation int64) T {
	var newView = NewViewP(leftX, rightX, topY, bottomY)
	return newRoot(newView, leafAllocation)
}

// A point with a slice of stored elements
type vpoint struct {
	x, y  float64
	elems []interface{}
}

// Indicates whether a vpoint is zeroed, i.e. uninitialised
func (np *vpoint) zeroed() bool {
	return np.elems == nil
}

// Resets this vpoint back to its uninitialised state
func (np *vpoint) zeroOut() {
	np.elems = nil
}

// Indicates whether a vpoint has the same x,y coords as those passed in
func (np *vpoint) sameLoc(x, y float64) bool {
	return np.x == x && np.y == y
}

func (np *vpoint) String() string {
	return fmt.Sprintf("(%v,%.3f,%.3f)", np.elems, np.x, np.y)
}

const LEAF_SIZE = 16

// A leaf struct implements the interface subtree. Like a node (see below),
// a leaf contains a View defining the rectangular area in which each vpoint
// could legally be located. A leaf struct may contain up to LEAF_SIZE non-zeroed
// vpoints.
// If any vpoint is non-empty, return false to zeroed(), then all vpoints 
// of lesser index are also non-empty i.e. if ps[3] is non-empty then so
// are ps[2], ps[1] and ps[0], while ps[4] or greater have no such constraints.
// The vpoints are not ordered in any way with respect to their geometric locations.
// A leaf is disposable if it was allocated outside the static leaf array, see root 
// below. If a leaf is marked as disposable it will not be recycled, but abandoned to
// the whimsy of the garbage collector.
type leaf struct {
	nextFree   *leaf
	view       View
	ps         [LEAF_SIZE]vpoint
	disposable bool
}

// Inserts each of the elements in elems into this leaf. There are three
// NB: We don't check that (x,y) is contained by this leaf's view, we rely of the parent
// node to ensure this.
// cases.
// If:			We find a non-empty vpoint at the exact location (x,y) 
//					- Append elems to this vpoint
// Else-If:		We find an empty vpoint available 
//					- Append elems to this vpoint and set the vpoint's location to (x,y)
// Else:		This leaf has overflowed 
//					- Replace this leaf with an intermediate node and re-allocate 
//					all of the elements in this leaf as well as those in elems into
//					the new node
func (l *leaf) insert(x, y float64, elems []interface{}, inPtr *subtree, r *root) {
	for i := range l.ps {
		if l.ps[i].zeroed() {
			l.ps[i].x = x
			l.ps[i].y = y
			l.ps[i].elems = elems
			return
		}
		if l.ps[i].sameLoc(x, y) {
			l.ps[i].elems = append(l.ps[i].elems, elems...)
			return
		}
	}
	// This leaf is full we need to create an intermediary node to divide it up
	newIntNode(x, y, elems, inPtr, l, r)
}

// This function creates a new node and adds all of the elements contained in l to it, 
// plus the new elements in elems. The pointer which previously pointed to l is 
// pointed at the new node. l is recycled.
func newIntNode(x, y float64, elems []interface{}, inPtr *subtree, l *leaf, r *root) {
	var newNode subtree
	newNode = r.newNode(l.View())
	for _, p := range l.ps {
		newNode.insert(p.x, p.y, p.elems, nil, r) // Does not require an inPtr param as we are passing into a *node
	}
	newNode.insert(x, y, elems, nil, r) // Does not require an inPtr param as we are passing into a *node
	*inPtr = newNode                    // Redirect the old leaf's reference to this intermediate node
	r.recycleLeaf(l)
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
// pred may have side-effects allowing for arbitrary processing of deleted elements.
// It is worth noting that if this leaf becomes empty it is the responsibility
// of this leaf's parent node to recycle it (when it chooses).
func (l *leaf) delete(view *View, pred func(x, y float64, e interface{}) bool, _ *subtree, _ *root) {
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

// Restores the leaf invariant that "if any vpoint is non-empty, then all vpoints 
// of lesser index are also non-empty" by rearranging the elements of ps.
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

// Returns a pointer to the View of this leaf
func (l *leaf) View() *View {
	return &l.view
}

// Sets the view for this leaf
func (l *leaf) setView(view *View) {
	l.view = *view
}

// Indicates whether or not this leaf contains any elements
func (l *leaf) isEmptyLeaf() bool {
	return l.ps[0].zeroed()
}

// Returns a human friendly string representation of this leaf
func (l *leaf) String() string {
	var str = l.view.String()
	for _, p := range l.ps {
		str += p.String()
	}
	return str
}

// A node struct implements the subtree interface.
// A node is the intermediate, non-leaf, storage structure for a 
// quadtree.
// It contains a View, indicating the rectangular area this node covers.
// Each subtree will have a view containing one of four quarters of
// this node's view. Every subtree is guaranteed to be non-nil and
// may be either a node or a leaf struct.
type node struct {
	nextFree   *node
	view       View
	children   [4]subtree
	disposable bool
}

// Inserts elems into the single child subtree whose view contains (x,y)
func (n *node) insert(x, y float64, elems []interface{}, _ *subtree, r *root) {
	for i := range n.children {
		if n.children[i].View().contains(x, y) {
			n.children[i].insert(x, y, elems, &n.children[i], r)
		}
	}
}

// Calls survey on each child subtree whose view overlaps with vs 
func (n *node) survey(vs []*View, fun func(x, y float64, e interface{})) {
	for i := range n.children {
		child := n.children[i]
		if overlaps(vs, child.View()) {
			n.children[i].survey(vs, fun)
		}
	}
}

// Calls delete on each child subtree whose view overlaps view
func (n *node) delete(view *View, pred func(x, y float64, e interface{}) bool, inPtr *subtree, r *root) {
	allEmpty := true
	for i := range n.children {
		if n.children[i].View().overlaps(view) {
			n.children[i].delete(view, pred, &n.children[i], r)
		}
		allEmpty = allEmpty && n.children[i].isEmptyLeaf()
	}
	if allEmpty && inPtr != nil {
		var l subtree
		l = r.newLeaf(n.View()) // TODO Think hard about whether this could error out
		*inPtr = l
		r.recycleNode(n)
	}
	return
}

// Returns the View for this node
func (n *node) View() *View {
	return &n.view
}

// Sets the view for this node
func (n *node) setView(view *View) {
	n.view = *view
}

// Always returns false - a node is never an empty leaf
func (n *node) isEmptyLeaf() bool {
	return false
}

// Returns a human friendly string representing this node, including its children.
func (n *node) String() string {
	return "<" + n.view.String() + "-\n" + n.children[0].String() + ", \n" + n.children[1].String() + ", \n" + n.children[2].String() + ", \n" + n.children[3].String() + ">"
}

// Each tree has a single root.
// The root is responsible for:
//	- Implementing the quadtree public interface T.
//	- Allocating and recycling leaf and node elements
type root struct {
	freeNode *node
	freeLeaf *leaf
	leaves   []leaf
	nodes    []node
	rootNode subtree
}

// Returns a new root ready for use as an empty quadtree
// Since for a quadtree with x leaves there will always be (x-1)/3 internal nodes
// we preallocate leafAllocation many leaves and (leafAllocation-1)/3 many nodes.
// 	NB: This number is not a hard limit, it only defines the number of statically allocated
// 	and managed nodes and leaves. More tree elements can be created and garbage will be garbage
// 	collected when they are recycled.
// A root node is initialised and the tree is ready for service.
func newRoot(view *View, leafAllocation int64) *root {
	if leafAllocation < 10 {
		leafAllocation = 10
	}
	leafNum := 3 - ((leafAllocation - 1) % 3) + leafAllocation
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
	rootNode := r.newNode(view)
	r.rootNode = rootNode
	return r
}

// Recursively recycle st and all of its children
func (r *root) recycle(st subtree) {
	switch st.(type) {
	case *leaf:
		r.recycleLeaf(st.(*leaf))
	case *node:
		r.recycleNode(st.(*node))
	}
}

// Returns a node with four leaves within the View provided.
// There are two kinds of node that can be returned.
// 	1: A free node from the roots static node array
//	2: A new node, marked disposable, fresh from the heap
// We only return 2 if 1 is not available.
func (r *root) newNode(view *View) (n *node) {
	if r.freeNode == nil {
		n = &node{view: *view, disposable: true}
	} else {
		n = r.freeNode
		r.freeNode = n.nextFree
		n.view = *view
	}
	r.newLeaves(view, &n.children)
	return
}

// Recycles n.
// If n is disposable this is a no-op, n should be garbage collected
// Otherwise n becomes r's next free node. r's old free node becomes 
// n's next free node.
// Each of n's children are recycled and n's children array is cleared.
// n's view is reset.
func (r *root) recycleNode(n *node) {
	if n.disposable {
		return
	}
	n.nextFree = r.freeNode
	r.freeNode = n
	for i := range n.children {
		r.recycle(n.children[i])
	}
	n.children = *new([4]subtree)
	n.view = *new(View)
}

// Returns a leaf with the view provided.
// There are two kinds of leaf that can be returned.
// 	1: A free leaf from the roots static leaf array
//	2: A new leaf, marked disposable, fresh from the heap
// We only return 2 if 1 is not available.
func (r *root) newLeaf(view *View) (l *leaf) {
	if r.freeLeaf == nil {
		l = &leaf{view: *view, disposable: true}
		return
	}
	l = r.freeLeaf
	r.freeLeaf = l.nextFree
	l.view = *view
	return
}

// Fills the array provided with new leaves each occupying 
// a quarter of the view provided.
func (r *root) newLeaves(view *View, leaves *[4]subtree) {
	v0, v1, v2, v3 := view.quarters()
	vs := []*View{v0, v1, v2, v3}
	for i := range leaves {
		leaves[i] = r.newLeaf(vs[i])
	}
	return
}

// Recycles l.
// If l is disposable this is a no-op, l should be garbage collected
// Otherwise, l becomes r's next free leaf. r's old free leaf becomes 
// l's next free leaf.
// l's view is reset. l's array of vpoints is reset.
func (r *root) recycleLeaf(l *leaf) {
	if l.disposable {
		return
	}
	l.nextFree = r.freeLeaf
	r.freeLeaf = l
	l.view = *new(View)
	l.ps = *new([LEAF_SIZE]vpoint)
}

// Inserts the value nval into this tree
func (r *root) Insert(x, y float64, nval interface{}) {
	elems := make([]interface{}, 1, 1)
	elems[0] = nval
	r.rootNode.insert(x, y, elems, nil, r)
}

// Deletes each element, e, under this node which satisfies two conditions
// 	1: e lies within view
//	2: pred(e) returns true
// pred may have side-effects allowing for arbitrary processing of deleted elements.
func (r *root) Delete(view *View, pred func(x, y float64, e interface{}) bool) {
	r.rootNode.delete(view, pred, nil, r)
}

// Applies fun to every element occurring within view in this tree
func (r *root) Survey(vs []*View, fun func(x, y float64, e interface{})) {
	r.rootNode.survey(vs, fun)
}

// Returns the View for this tree
func (r *root) View() *View {
	return r.rootNode.View()
}

// Counts the number of free nodes available.
// For debugging only
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

// Counts the number of free leaves available.
// For debugging only
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
