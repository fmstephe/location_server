package quadtree

// Public interface for quadtrees.
type T interface {
	View() *View
	// Inserts e into this quadtree at point (x,y)
	Insert(x, y float64, e interface{})
	// Applies fun to every element in this quadtree that lies within any view in views
	Survey(views []*View, fun func(x, y float64, e interface{}))
	// Applies pred to every element in this quadtree that lies within any view in views
	// If pred returns true that element is removed
	Del(views *View, pred func(x, y float64, e interface{}) bool)
	// Provides a human readable (as far as possible) string representation of this tree
	String() string
}
