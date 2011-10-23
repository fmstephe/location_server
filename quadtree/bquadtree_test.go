package quadtree

import (
	"testing"
	"container/vector"
	"strconv"
)

var width, height = 100.0, 100.0
var treeNum = 10
var pointsSmall = 1000
var pointsLarge = 10000
var repsSingle = 1
var repsLarge = 10

func makeTrees(tNum int, w, h float64) []T {
	trees := make([]T, tNum)
	for i := range trees {
		trees[i] = NewQuadTree(0, w, 0, h)
	}
	return trees
}

func makePoints(tNum, pNum int, w, h float64) [][]point {
	points := make([][]point, tNum)
	for i := range points {
		points[i] = fillView(OrigViewP(w, h), pNum)
	}
	return points
}

func makeFilledTrees(tNum, pNum, reps int, w, h float64) []T {
	trees := makeTrees(tNum, w, h)
	points := makePoints(tNum, pNum, w, h)
	for ti := range trees {
		tree := trees[ti]
		ps := points[ti]
		for pi, p := range ps {
			for ri := 0; ri < reps; ri++ {
				name := strconv.Itoa(pi) + "_test_" + strconv.Itoa(ri)
				tree.Insert(p.x, p.y, name)
			}
		}
	}
	return trees
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		trees := makeTrees(treeNum, width, height)
		treePoints := makePoints(treeNum, pointsLarge, width, height)
		b.StartTimer()
		for j := range trees {
			tree := trees[j]
			ps := treePoints[j]
			for _, p := range ps {
				tree.Insert(p.x, p.y, "test")
			}
		}
	}
}

func BenchmarkSurveyReps(b *testing.B) {
	benchmarkSurvey(b, pointsSmall, repsLarge)
}

func BenchmarkSurveySingle(b *testing.B) {
	benchmarkSurvey(b, pointsLarge, repsSingle)
}

func benchmarkSurvey(b *testing.B, pointNum, repNum int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		trees := makeFilledTrees(treeNum, pointNum, repNum, width, height)
		b.StartTimer()
		for j := range trees {
			tree := trees[j]
			collected := new(vector.Vector)
			collected.Resize(pointNum, pointNum)
			fun := func(x, y float64, e interface{}) {
				collected.Push(e)
			}
			tree.Survey([]*View{tree.View()}, fun)
		}
	}
}

func BenchmarkDeleteReps(b *testing.B) {
	benchmarkDelete(b, pointsSmall, repsLarge)
}

func BenchmarkDeleteSingle(b *testing.B) {
	benchmarkDelete(b, pointsLarge, repsSingle)
}

func benchmarkDelete(b *testing.B, pointNum, repNum int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		trees := makeFilledTrees(treeNum, pointNum, repNum, width, height)
		b.StartTimer()
		for j := range trees {
			tree := trees[j]
			deleted := new(vector.Vector)
			deleted.Resize(pointNum, pointNum)
			q1, q2, q3, q4 := tree.View().quarters()
			pred := func(x, y float64, e interface{}) bool {
				deleted.Push(e)
				return true
			}
			tree.Delete(q1, pred)
			tree.Delete(q2, pred)
			tree.Delete(q3, pred)
			tree.Delete(q4, pred)
		}
	}
}
