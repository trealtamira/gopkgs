package tiling

import (
	"math"
)

//Tile reresent a tile in ZoomLevel
type Tile struct {
	X int
	Y int
	Z int
}

//ZoomLevel represent a single zoom level of the tile map pyramidal system
type ZoomLevel struct {
	zoom    int
	mxSize  float64
	hLength float64
	vLength float64
}

//NewZoomLevel create a new zoomlevel instance at level z
func NewZoomLevel(z int) *ZoomLevel {
	matrixSize := float64(math.Pow(2, float64(z)))
	hTileLength := equator / matrixSize
	vTileLength := (2 * meridian) / matrixSize
	zl := ZoomLevel{zoom: z, mxSize: matrixSize, hLength: hTileLength, vLength: vTileLength}
	return &zl
}

//Cardinality gives the number of tiles in the zoom level
func (z *ZoomLevel) Cardinality() float64 {
	return math.Pow(z.mxSize, 2)
}

//TileFor gives the tile coordinates for the given point for the current zoom level
func (z *ZoomLevel) TileFor(m PointM) Tile {
	x := math.Floor((m.E + (equator / 2)) / z.hLength)
	y := math.Floor((meridian - m.N) / z.vLength)
	t := Tile{X: int(x), Y: int(y), Z: int(z.zoom)}
	return t
}

//Extent return the extent mercator coordinates of the given tile
func (z *ZoomLevel) Extent(x, y int) ExtentM {
	minEast := (float64(x) * z.hLength) - (equator / 2)
	maxEast := (float64(x+1) * z.hLength) - (equator / 2)
	maxNorth := meridian - (float64(y) * z.vLength)
	minNorth := meridian - (float64(y+1) * z.vLength)
	ne := PointM{N: maxNorth, E: maxEast}
	sw := PointM{N: minNorth, E: minEast}
	me := ExtentM{NE: ne, SW: sw}
	return me
}

//Extent return the mercator extent coordinates of the given tile
func (t *Tile) Extent() ExtentM {
	z := NewZoomLevel(t.Z)
	return z.Extent(t.X, t.Y)
}
