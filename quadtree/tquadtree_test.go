package quadtree

import (
	"container/list"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type dim struct {
	width, height float64
}

const dups = 10
const treeLim = 10000

var testTrees []T

type point struct {
	x, y float64
}

var testRand = rand.New(rand.NewSource(time.Nanoseconds()))

func init() {
	testTrees = []T{
		NewQuadTree(0, 10, 0, 10, treeLim),
		NewQuadTree(0, 1, 0, 2, treeLim),
		NewQuadTree(0, 100, 0, 300, treeLim),
		NewQuadTree(0, 20.4, 0, 35.6, treeLim),
		NewQuadTree(0, 1e10, 0, 500.00000001, treeLim),
		// Negative regions
		NewQuadTree(-10, 10, -10, 10, treeLim),
		NewQuadTree(-1, 1, -2, 2, treeLim),
		NewQuadTree(-100, 100, -300, 300, treeLim),
		NewQuadTree(-20.4, 20.4, -35.6, 35.6, treeLim),
		NewQuadTree(-1e10, 1e10, -500.00000001, 500.00000001, treeLim),
	}
}

func clearTrees() {
	for i := range testTrees {
		testTrees[i].Delete(testTrees[i].View(), SimpleDelete())
	}
}

// Test that we can insert a single element into the tree and then retrieve it
func TestOneElement(t *testing.T) {
	clearTrees()
	for _, tree := range testTrees {
		testOneElement(tree, t)
	}
}

func testOneElement(tree T, t *testing.T) {
	x, y := randomPosition(tree.View())
	tree.Insert(x, y, "test")
	fun, results := SimpleSurvey()
	tree.Survey([]*View{tree.View()}, fun)
	if results.Len() != 1 || "test" != results.Front().Value {
		t.Errorf("Failed to find required element at (%f,%f), in tree \n%v", x, y, tree)
	}
}

// Test that if we add 5 elements into a single quadrant of a fresh tree
// We can successfully retrieve those elements. This test is tied to
// the implementation detail that a quadrant with 5 elements will 
// over-load a single leaf and must rearrange itself to fit the 5th 
// element in.
func TestFullLeaf(t *testing.T) {
	clearTrees()
	for _, tree := range testTrees {
		v1, v2, v3, v4 := tree.View().quarters()
		testFullLeaf(tree, v1, "v1", t)
		testFullLeaf(tree, v2, "v2", t)
		testFullLeaf(tree, v3, "v3", t)
		testFullLeaf(tree, v4, "v4", t)
	}
}

func testFullLeaf(tree T, v *View, msg string, t *testing.T) {
	for i := 0; i < 5; i++ {
		x, y := randomPosition(v)
		name := "test" + strconv.Itoa(i)
		tree.Insert(x, y, name)
	}
	fun, results := SimpleSurvey()
	tree.Survey([]*View{v}, fun)
	if results.Len() != 5 {
		t.Error(msg, "Inserted 5 elements into a fresh quadtree and retrieved only ", results.Len())
	}
}

// Tests that we can add a large number of random elements to a tree
// and create random views for collecting from the populated tree.
func TestScatter(t *testing.T) {
	clearTrees()
	for _, tree := range testTrees {
		testScatter(tree, t)
	}
	clearTrees()
	for _, tree := range testTrees {
		testScatterDup(tree, t)
	}
}

func testScatter(tree T, t *testing.T) {
	ps := fillView(tree.View(), 1000)
	for _, p := range ps {
		tree.Insert(p.x, p.y, "test")
	}
	for i := 0; i < 1000; i++ {
		sv := subView(tree.View())
		var count int
		for _, v := range ps {
			if sv.contains(v.x, v.y) {
				count++
			}
		}
		fun, results := SimpleSurvey()
		tree.Survey([]*View{sv}, fun)
		if count != results.Len() {
			t.Errorf("Failed to retrieve %d elements in scatter test, found only %d", count, results.Len())
		}
	}
}

// Tests that we can add multiple elements to the same location
// and still retrieve all elements, including duplicates, using 
// randomly generated views.
func testScatterDup(tree T, t *testing.T) {
	ps := fillView(tree.View(), 1000)
	for _, p := range ps {
		for i := 0; i < dups; i++ {
			tree.Insert(p.x, p.y, "test_"+strconv.Itoa(i))
		}
	}
	for i := 0; i < 1000; i++ {
		sv := subView(tree.View())
		var count int
		for _, v := range ps {
			if sv.contains(v.x, v.y) {
				count++
			}
		}
		fun, results := SimpleSurvey()
		tree.Survey([]*View{sv}, fun)
		if count*dups != results.Len() {
			t.Error("Failed to retrieve %i elements in duplicate scatter test, found only %i", count*dups, results.Len())
		}
	}
}

// Tests that when we
// 1: Add a single element to an empty tree
// 2: Remove that element from the tree
// We get
// 1: The single element is the only element in the deleted list
// 2: The tree no longer contains any elements
func TestSimpleAddDelete(t *testing.T) {
	clearTrees()
	for _, tree := range testTrees {
		testAddDelete(tree, t)
	}
	clearTrees()
	for _, tree := range testTrees {
		testAddDeleteDup(tree, t)
	}
	clearTrees()
	for _, tree := range testTrees {
		testAddDeleteMulti(tree, t)
	}
}

// Add element delete everything from the tree.
func testAddDelete(tree T, t *testing.T) {
	elem := "test"
	//testDeleteSimple(NewQuadTree(0, d.width, 0, d.height, 10000), []interface{}{elem}, []interface{}{elem}, false, "Simple Global Delete", t)
	testDeleteSimple(tree, []interface{}{elem}, []interface{}{elem}, true, "Simple Exact Delete", t)
}

// Add two elements, delete one element from entire tree
func testAddDeleteDup(tree T, t *testing.T) {
	elem := "test"
	elemII := "testII"
	//testDeleteSimple(NewQuadTree(0, d.width, 0, d.height, 10000), []interface{}{elem, elemII}, []interface{}{elem}, false, "Simple Gobal Delete Take One Of Two", t)
	testDeleteSimple(tree, []interface{}{elem, elemII}, []interface{}{elem}, false, "Simple Exact Delete Take One Of Two", t)
}

// Add two elements, delete both from entire tree
func testAddDeleteMulti(tree T, t *testing.T) {
	elem := "test"
	elemII := "testII"
	//testDeleteSimple(NewQuadTree(0, d.width, 0, d.height, 10000), []interface{}{elem, elemII}, []interface{}{elem, elemII}, false, "Simple Global Delete Take Two Of Two", t)
	testDeleteSimple(tree, []interface{}{elem, elemII}, []interface{}{elem, elemII}, true, "Simple Exact Delete Take Two Of Two", t)
}

// Tests a very limited deletion scenario. Here we will insert every element in 'insert' into the tree at a
// single random point. Then we will delete every element in delete from the tree. 
// If exact == true then the view used to delete covers eactly the insertion point. Otherwise, it covers the
// entire tree.
// We assert that every element of delete has been deleted from the tree (testDelete)
// We assert that every element in insert but not in delete is still in the tree (testSurvey)
// errPrfx is used to distinguish the error messages from different tests using this method.
func testDeleteSimple(tree T, insert, delete []interface{}, exact bool, errPrfx string, t *testing.T) {
	x, y := randomPosition(tree.View())
	for _, e := range insert {
		tree.Insert(x, y, e)
	}
	expCol := new(list.List)
OUTER_LOOP:
	for _, i := range insert {
		for _, d := range delete {
			if i == d {
				continue OUTER_LOOP
			}
		}
		expCol.PushBack(i)
	}
	expDel := new(list.List)
	for _, d := range delete {
		expDel.PushBack(d)
	}
	pred, deleted := makeDelClosure(delete)
	delView := tree.View()
	if exact {
		delView = NewViewP(x, x, y, y)
	}
	testDelete(tree, delView, pred, deleted, expDel, t, errPrfx)
	fun, collected := SimpleSurvey()
	testSurvey(tree, tree.View(), fun, collected, expCol, t, errPrfx)
}

func TestScatterDelete(t *testing.T) {
	clearTrees()
	for _, tree := range testTrees {
		testScatterDelete(tree, t)
	}
	clearTrees()
	for _, tree := range testTrees {
		testScatterDeleteMulti(tree, t)
	}
}

func testScatterDelete(tree T, t *testing.T) {
	name := "test"
	pointNum := 1000
	ps := fillView(tree.View(), pointNum)
	for i, p := range ps {
		tree.Insert(p.x, p.y, name+strconv.Itoa(i))
	}
	delView := subView(tree.View())
	expDel := new(list.List)
	expCol := new(list.List)
	for i, p := range ps {
		if delView.contains(p.x, p.y) {
			expDel.PushBack(name + strconv.Itoa(i))
		} else {
			expCol.PushBack(name + strconv.Itoa(i))
		}
	}
	pred, deleted := CollectingDelete()
	testDelete(tree, delView, pred, deleted, expDel, t, "Scatter Insert and Delete Under Area")
	fun, collected := SimpleSurvey()
	testSurvey(tree, tree.View(), fun, collected, expCol, t, "Scatter Insert and Delete Under Area")
}

func testScatterDeleteMulti(tree T, t *testing.T) {
	name := "test"
	pointNum := 1000
	points := fillView(tree.View(), pointNum)
	for i, p := range points {
		for d := 0; d < dups; d++ {
			tree.Insert(p.x, p.y, name+strconv.Itoa(i)+"_"+strconv.Itoa(d))
		}
	}
	delView := subView(tree.View())
	expDel := new(list.List)
	expCol := new(list.List)
	for i, p := range points {
		if delView.contains(p.x, p.y) {
			for d := 0; d < dups; d++ {
				expDel.PushBack(name + strconv.Itoa(i) + "_" + strconv.Itoa(d))
			}
		} else {
			for d := 0; d < dups; d++ {
				expCol.PushBack(name + strconv.Itoa(i) + "_" + strconv.Itoa(d))
			}
		}
	}
	pred, deleted := CollectingDelete()
	testDelete(tree, delView, pred, deleted, expDel, t, "Scatter Insert and Delete Under Area With Three Elements Per Location")
	fun, results := SimpleSurvey()
	testSurvey(tree, tree.View(), fun, results, expCol, t, "Scatter Insert and Delete Under Area With Three Elements Per Location")
}

// Test that we can successfully insert substantially more elements into
// the tree than the minLeafMax value passed into NewQuadTree(...)
// It is important to note that trees can usually absorb between 6-9 times minLeafMax
func TestOverloadLifecycle(t *testing.T) {
	tree := NewQuadTree(0, 10, 0, 10, 10)
	points := fillView(tree.View(), 100000)
	for _, p := range points {
		tree.Insert(p.x, p.y, "overload")
	}
}

// Handy function to print out at the expected number of random elements we can add to a tree
// for a given tree lim size - curiously wavers between 6 - 9 in a curious wave pattern
/*
func disabledTestInsertLimits(t *testing.T) {
	for treeLim := 10000; treeLim < 1000000; treeLim += 10000 {
		pointNum := treeLim * 10
		tree := NewQuadTree(0, 10, 0, 10, int64(treeLim))
		points := fillView(tree.View(), pointNum)
		for pi, p := range points {
			if err := tree.Insert(p.x, p.y, "limit"); err != nil {
				fmt.Println(float64(pi) / float64(treeLim))
				break
			}
		}
	}
}
*/

func testDelete(tree T, view *View, pred func(x, y float64, e interface{}) bool, deleted, expDel *list.List, t *testing.T, errPfx string) {
	tree.Delete(view, pred)
	if deleted.Len() != expDel.Len() {
		t.Errorf("%s: Expecting %v deleted element(s), found %v", errPfx, expDel.Len(), deleted.Len())
	}
OUTER_LOOP:
	for i := expDel.Front(); i != nil; i = i.Next() {
		for j := deleted.Front(); j != nil; j = j.Next() {
			expVal := i.Value
			delVal := j.Value
			if expVal == delVal {
				continue OUTER_LOOP
			}
		}
		t.Errorf("%s: Expecting to find %v in deleted vector, was not found", errPfx, i.Value)
	}
}

func testSurvey(tree T, view *View, fun func(x, y float64, e interface{}), collected, expCol *list.List, t *testing.T, errPfx string) {
	tree.Survey([]*View{view}, fun)
	if collected.Len() != expCol.Len() {
		t.Errorf("%s: Expecting %v collected element(s), found %v", errPfx, expCol.Len(), collected.Len())
	}
	/* This code checks that every expected element is present
		   In practice this is too slow - disabled
	OUTER_LOOP:
		for i := 0; i < expCol.Len(); i++ {
			expVal := expCol.At(i)
			for j := 0; j < collected.Len(); j++ {
				colVal := collected.At(j)
				if expVal == colVal {
					continue OUTER_LOOP
				}
			}
			t.Errorf("%s: Expecting to find %v in collected vector, was not found", errPfx, expCol.At(i))
		}
	*/
}

// Creates a closure which deletes all elements which are present in elem
// Returns the closure plus a list.List into which deleted elements are accumulated
func makeDelClosure(elems []interface{}) (pred func(x, y float64, e interface{}) bool, deleted *list.List) {
	deleted = new(list.List)
	pred = func(x, y float64, e interface{}) bool {
		for i := range elems {
			if e == elems[i] {
				deleted.PushBack(e)
				return true
			}
		}
		return false
	}
	return
}

func randomPosition(v *View) (x, y float64) {
	x = testRand.Float64()*(v.rx-v.lx) + v.lx
	y = testRand.Float64()*(v.by-v.ty) + v.ty
	return
}

func fillView(v *View, c int) []point {
	ps := make([]point, c)
	for i := 0; i < c; i++ {
		x, y := randomPosition(v)
		ps[i] = point{x: x, y: y}
	}
	return ps
}

func subView(v *View) *View {
	lx := testRand.Float64()*(v.rx-v.lx) + v.lx
	rx := testRand.Float64()*(v.rx-lx) + lx
	ty := testRand.Float64()*(v.by-v.ty) + v.ty
	by := testRand.Float64()*(v.by-ty) + ty
	return NewViewP(lx, rx, ty, by)
}
