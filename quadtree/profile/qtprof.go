package main

import (
	"flag"
	"location_server/loc_server/quadtree"
	"math/rand"
	"os"
	"runtime/pprof"
)

const iterations = 1
const elemCount = 1000000
const treeSize = 100

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, _ := os.Create(*cpuprofile)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	for i := 0; i < iterations; i++ {
		tree := quadtree.NewQuadTree(0, treeSize, 0, treeSize, elemCount/6)
		for i := 0; i < elemCount; i++ {
			x := rand.Float64() * treeSize
			y := rand.Float64() * treeSize
			tree.Insert(x, y, i)
		}
		vs := []*quadtree.View{tree.View()}
		col := make([]interface{}, 0, elemCount)
		fun := func(x, y float64, e interface{}) {
			col = append(col, e)
		}
		tree.Survey(vs, fun)
		tree.Delete(tree.View(), quadtree.SimpleDelete())
		println(col)
	}
}
