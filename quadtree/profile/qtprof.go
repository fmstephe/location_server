package main

import (
	"flag"
	"github.com/fmstephe/flib/fstrconv"
	"github.com/fmstephe/location_server/quadtree"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

const iterations = 1
const elemCount = 1000000
const treeSize = 1000

var cpuprofile = flag.String("file", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, _ := os.Create(*cpuprofile)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	start := time.Now()
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
		tree.Del(tree.View(), quadtree.SimpleDelete())
		println(col)
	}
	total := time.Now().Sub(start)
	println(fstrconv.ItoaComma(elemCount))
	println(total.String())
	secs := total.Nanoseconds() / (1000 * 1000 * 1000)
	if secs == 0 {
		return
	}
	println(secs)
	perSec := elemCount / secs
	println(fstrconv.ItoaComma(perSec), "elements per second")
}
