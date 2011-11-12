package locserver

import (
	"math"
)

const (
	// Useful constants for converting between radians and degrees
	RadToDeg = 180 / math.Pi
	DegToRad = math.Pi / 180

	// Distance, in meters covered across one lattitudinal degree at the equator (at sealevel blah, blah, blah)
	// Lattitude distance is constant because we make the unsophisticated assumption that the world is spherical
	// NB: Under the WGS84 reference spheroid metresPerLat should vary between 110,574 (at the equator) and 111,468 (at the poles)
	metresPerLat = 110574.0

	// Distance, in meters, covered across one longitudinal degree at the equator (at sealevel blah, blah, blah)
	// Distance covered changes across different lattitudes as the longitundial meridians converge at the poles.
	// According to wikipedia this relationship is (pi/180)cos(lat)M, where M is the radius of the earth (avg ~6,367,449 metres)
	metresPerLng = 111312.0

	// The mean radius of the earth in metres
	metresEarthRadius = 6367449.0
	metresEarthCirc   = 40075160.0

	maxNorthDeg = 90
	maxSouthDeg = -90
	maxEastDeg  = 180
	maxWestDeg  = -180

	maxNorthMetres = maxNorthDeg * metresPerLat
	maxSouthMetres = maxSouthDeg * metresPerLat
	maxEastMetres  = maxEastDeg * metresPerLng
	maxWestMetres  = maxWestDeg * metresPerLng
)

// Takes a (lat,lng) point and returns a (mLat,mLng) point
// where:
// 		mNS is the number of north/south metres from (lat,lng) (0,0)
//			A positive indicates a position north of (0,0)
//			A negative indicates a position south of (0,0)
// 		mEW is the number of east/west metres from (0,0)
//			A positive indicates a position east of (0,0)
//			A negative indicates a position west of (0,0)
func metresFromOrigin(lat, lng float64) (mNS, mEW float64) {
	lngRat := math.Abs(lat/maxNorthDeg)
	widthRat := math.Sqrt(1 - lngRat)
	mEW = lng * metresPerLng * widthRat
	mNS = lat * metresPerLat
	return
}
