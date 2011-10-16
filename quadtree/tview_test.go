package quadtree

import (
	"testing"
	"rand"
	"time"
	"math"
)

// Width and Height struct
type wh struct {
	width, height float64
}

// Slice of width/height pairs for testing origin views
var whs = []wh{
	wh{10.0, 10.0},
	wh{5.5, 30.03},
	wh{0.0, 0.0},
	wh{9999.9999, 1.123456},
	wh{10.01e8, 4.4e3},
}

// Struct of four points
type fp struct {
	lx, rx, ty, by float64
}

// Slice of four point structs for testing literal views
var fps = []fp{
	fp{10.0, 10.0, 10.0, 10.0},
	fp{10.0, 10.0, 10.0, 10.0},
	fp{5.5, 30.03, 3.45, 5.96},
	fp{0.0, 0.0, 0.0, 0.0},
	fp{1.123456, 9999.9999, 9876.5, 12345.5},
	// Negative ones
	fp{-4.4e3, 10.01e8, -45.0e4, 5.0e5},
	fp{-5.5, 30.03, -3.45, 5.96},
	fp{-0.0, 0.0, -0.0, 0.0},
	fp{-1.123456, 9999.9999, -9876.5, 12345.5},
	fp{-4.4e3, 10.01e8, -45.0e4, 5.0e5},
}

// Slice of illegal four point structs for testing literal views
var ifps = []fp{
	fp{5.5, 5.4, 3.45, 5.96},
	fp{5.5, 5.6, 3.45, 3.44},
	fp{5.5, 5.4, 3.45, 3.44},
}

func TestMallocView(t *testing.T) {
	v := new(View)
	if v.lx != 0 {
		t.Error("Left x not zeroed")
	}
	if v.rx != 0 {
		t.Error("Right x not zeroed")
	}
	if v.ty != 0 {
		t.Error("Top y not zeroed")
	}
	if v.by != 0 {
		t.Error("Bottom y not zeroed")
	}
}

func TestOrigView(t *testing.T) {
	for _, v := range whs {
		testOrigView(v.width, v.height, t)
	}
}

func testOrigView(width, height float64, t *testing.T) {
	v := OrigView(width, height)
	if v.lx != 0 {
		t.Error("Left x not at origin")
	}
	if v.rx != width {
		t.Errorf("Right x %10.3f : expecting %10.3f", v.rx, width)
	}
	if v.ty != 0 {
		t.Error("Top y not at origin")
	}
	if v.by != height {
		t.Errorf("Bottom y %10.3f : expecting %10.3f", v.by, height)
	}
}

func TestNewView(t *testing.T) {
	for _, v := range fps {
		testNewView(v.lx, v.rx, v.ty, v.by, t)
	}
}

func testNewView(lx, rx, ty, by float64, t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Test panic, failure for %10.3f, %10.3f, %10.3f, %10.3f", lx, rx, ty, by)
		}
	}()
	v := NewView(lx, rx, ty, by)
	if v.lx != lx {
		t.Errorf("Left x %10.3f : expecting %10.3f", v.lx, lx)
	}
	if v.rx != rx {
		t.Errorf("Right x %10.3f : expecting %10.3f", v.rx, rx)
	}
	if v.ty != ty {
		t.Errorf("Right x %10.3f : expecting %10.3f", v.ty, ty)
	}
	if v.by != by {
		t.Errorf("Bottom y %10.3f : expecting %10.3f", v.by, by)
	}
}

func TestIllegalView(t *testing.T) {
	for _, v := range ifps {
		testIllegalView(v.lx, v.rx, v.ty, v.by, t)
	}
}

func testIllegalView(lx, rx, ty, by float64, t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Test did not panic, %10.3f, %10.3f, %10.3f, %10.3f should be illegal", lx, rx, ty, by)
		}
	}()
	NewView(lx, rx, ty, by)
}

func TestOverLap(t *testing.T) {
	rand.Seed(time.Nanoseconds())
	for i := 0; i < 10000; i++ {
		v1, v2 := overlap()
		if !v1.overlaps(v2) {
			t.Errorf("View %v and %v not reported as overlapping", v1, v2)
			t.Error("+------------------------------------------------+")
		}
		if !v2.overlaps(v1) {
			t.Errorf("View %v and %v not reported as overlapping", v2, v1)
			t.Error("<------------------------------------------------>")
		}
	}
}

func TestDisjoint(t *testing.T) {
	rand.Seed(time.Nanoseconds())
	for i := 0; i < 1000; i++ {
		v1, v2 := disjoint()
		if v1.overlaps(v2) {
			t.Errorf("View %v and %v are reported as overlapping", v1, v2)
			t.Error("+------------------------------------------------+")
		}
		if v2.overlaps(v1) {
			t.Errorf("View %v and %v are reported as overlapping", v2, v1)
			t.Error("<------------------------------------------------>")
		}
	}
}

func TestSubtract(t *testing.T) {
	var v1, v2 *View
	var eqv []*View
	// Left side
	v1 = NewViewP(0, 2, 0, 2)
	v2 = NewViewP(1, 2, 0, 2)
	eqv = []*View{NewViewP(0, 1, 0, 2)}
	testSubtract(t, v1, v2, eqv, 1, "Left side")
	// Right Side
	v1 = NewViewP(0, 2, 0, 2)
	v2 = NewViewP(0, 1, 0, 2)
	eqv = []*View{NewViewP(1, 2, 0, 2)}
	testSubtract(t, v1, v2, eqv, 1, "Right side")
	// Bottom Side
	v1 = NewViewP(0, 2, 0, 2)
	v2 = NewViewP(0, 2, 1, 2)
	eqv = []*View{NewViewP(0, 2, 0, 1)}
	testSubtract(t, v1, v2, eqv, 1, "Bottom side")
	// Top Side
	v1 = NewViewP(0, 2, 0, 2)
	v2 = NewViewP(0, 2, 0, 1)
	eqv = []*View{NewViewP(0, 2, 1, 2)}
	testSubtract(t, v1, v2, eqv, 1, "Top side")
	// Bottom Left
	v1 = NewViewP(0, 2, 0, 2)
	v2 = NewViewP(1, 2, 0, 1)
	eqv = []*View{NewViewP(0, 1, 0, 2), NewViewP(0, 2, 1, 2)}
	testSubtract(t, v1, v2, eqv, 2, "Bottom left")
	// Bottom Right
	v1 = NewViewP(0, 2, 0, 2)
	v2 = NewViewP(0, 1, 0, 1)
	eqv = []*View{NewViewP(1, 2, 0, 2), NewViewP(0, 2, 1, 2)}
	testSubtract(t, v1, v2, eqv, 2, "Bottom right")
	// Top Left
	v1 = NewViewP(0, 2, 0, 2)
	v2 = NewViewP(1, 2, 1, 2)
	eqv = []*View{NewViewP(0, 1, 0, 2), NewViewP(0, 2, 0, 1)}
	testSubtract(t, v1, v2, eqv, 2, "Top left")
	// Top Right
	v1 = NewViewP(0, 2, 0, 2)
	v2 = NewViewP(0, 1, 1, 2)
	eqv = []*View{NewViewP(1, 2, 0, 2), NewViewP(0, 2, 0, 1)}
	testSubtract(t, v1, v2, eqv, 2, "Top left")
	// Left Right
	v1 = NewViewP(0, 3, 0, 3)
	v2 = NewViewP(1, 2, -20, 20)
	eqv = []*View{NewViewP(0, 1, 0, 3), NewViewP(2, 3, 0, 3)}
	testSubtract(t, v1, v2, eqv, 2, "Left right")
	// Top Bottom
	v1 = NewViewP(0, 3, 0, 3)
	v2 = NewViewP(-20, 20, 1, 2)
	eqv = []*View{NewViewP(0, 3, 0, 1), NewViewP(0, 3, 2, 3)}
	testSubtract(t, v1, v2, eqv, 2, "Top bottom")
	// Centre
	v1 = NewViewP(0, 4, 0, 4)
	v2 = NewViewP(1, 2, 1, 2)
	eqv = []*View{NewViewP(0, 1, 0, 4), NewViewP(2, 4, 0, 4), NewViewP(0, 4, 0, 1), NewViewP(0, 4, 2, 4)}
	testSubtract(t, v1, v2, eqv, 4, "Centre")
}

func testSubtract(t *testing.T, v1, v2 *View, eqv []*View, vNum int, prefix string) {
	vs := v1.Subtract(v2)
	if len(vs) != vNum {
		t.Errorf("%s subtract, expecting %d view found %d", prefix, vNum, len(vs))
	}
	for ei := range eqv {
		for vi := range vs {
			if eqv[ei].eq(vs[vi]) {
				break
			}
			if vi == len(vs)-1 {
				t.Errorf("%s subtract, expected %v found %v", prefix, eqv, vs)
				return
			}
		}
	}
}

func overlap() (v1, v2 *View) {
	lx, rx := oPair(negRFLoat64(), negRFLoat64())
	ty, by := oPair(negRFLoat64(), negRFLoat64())
	x1 := rand.Float64()
	y1 := rand.Float64()
	nx, ny := nearest(x1, y1, lx, rx, ty, by)
	var x2, y2 float64
	if x1 > nx {
		x2 = nx - rand.Float64()
	} else {
		x2 = rand.Float64() + nx
	}
	if y1 > ny {
		y2 = ny - rand.Float64()
	} else {
		y2 = rand.Float64() + ny
	}
	lx2, rx2 := oPair(x1, x2)
	ty2, by2 := oPair(y1, y2)
	return NewViewP(lx, rx, ty, by), NewViewP(lx2, rx2, ty2, by2)
}

func disjoint() (v1, v2 *View) {
	lx, rx := oPair(negRFLoat64(), negRFLoat64())
	ty, by := oPair(negRFLoat64(), negRFLoat64())
	v1 = NewViewP(lx, rx, ty, by)
	var x1, y1 float64
	for true {
		x1 = negRFLoat64()
		y1 = negRFLoat64()
		if !v1.contains(x1, y1) {
			break
		}
	}
	nx, ny := nearest(x1, y1, lx, rx, ty, by)
	var x2, y2 float64
	if x1 < nx {
		x2 = nx - rand.Float64()
	} else {
		x2 = rand.Float64() + nx
	}
	if y1 < ny {
		y2 = ny - rand.Float64()
	} else {
		y2 = rand.Float64() + ny
	}
	lx2, rx2 := oPair(x1, x2)
	ty2, by2 := oPair(y1, y2)
	v2 = NewViewP(lx2, rx2, ty2, by2)
	return
}

func oPair(f1, f2 float64) (r1, r2 float64) {
	r1 = math.Fmin(f1, f2)
	r2 = math.Fmax(f1, f2)
	return
}

func nearest(x, y, lx, rx, ty, by float64) (nx, ny float64) {
	d1 := dist(x, y, lx, ty)
	d2 := dist(x, y, rx, ty)
	d3 := dist(x, y, lx, by)
	d4 := dist(x, y, rx, by)
	n1 := math.Fmin(d1, d2)
	n2 := math.Fmin(n1, d3)
	n3 := math.Fmin(n2, d4)
	if n3 == d1 {
		return lx, ty
	}
	if n3 == d2 {
		return rx, ty
	}
	if n3 == d3 {
		return lx, by
	}
	if n3 == d4 {
		return rx, by
	}
	panic("Shouldn't reach this line")
}

func dist(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x1-x2, 2.0) + math.Pow(y1-y2, 2.0))
}

func negRFLoat64() float64 {
	f := rand.Float64()
	d := rand.Float64()
	if d < 0.5 {
		return -f
	}
	return f
}
