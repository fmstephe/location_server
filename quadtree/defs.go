package quadtree

// Public interface for quadtrees.
// Note that only root implements this interface, not leaf.
type T interface {
	View() *View
	// Error may be returned to indicate that e could not be inserted if the 
	// implementation has size restrictions
	Insert(x, y float64, e interface{})
	//
	Survey(view []*View, fun func(x, y float64, e interface{}))
	//
	Delete(view *View, pred func(x, y float64, e interface{}) bool)
	//
	String() string
}
