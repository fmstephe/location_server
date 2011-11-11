package quadtree

import(
	"container/list"
)

//
//	A Simple quadtree collector which will push every element into col
//
func SimpleSurvey() (fun func(x, y float64, e interface{}), col *list.List) {
	col = list.New()
	fun = func(x, y float64, e interface{}) {
		col.PushBack(e)
	}
	return
}

//
//	A Simple quadtree delete function which indicates that every element given to it should be deleted
//
func SimpleDelete() (pred func(x, y float64, e interface{}) bool) {
	pred = func(x, y float64, e interface{}) bool {
		return true
	}
	return
}

//
//	A quadtree delete function which indicates that every element given to it should be deleted.
//	Additionally each element deleted will be pushed into col
//
func CollectingDelete() (pred func(x, y float64, e interface{}) bool, col *list.List) {
	col = list.New()
	pred = func(x, y float64, e interface{}) bool {
		col.PushBack(e)
		return true
	}
	return
}

// 
//	Determines if a point lies inside at least one of a slice of *View
//
func contains(vs []*View, x, y float64) bool {
	for _, v := range vs {
		if v.contains(x, y) {
			return true
		}
	}
	return false
}

//
//	Determines if a view overlaps at least one of a slice of *View
//
func overlaps(vs []*View, oV *View) bool {
	for _, v := range vs {
		if oV.overlaps(v) {
			return true
		}
	}
	return false
}
