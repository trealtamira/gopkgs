package tiling

import (
	"fmt"
	"math"
)

const (
	tileMaxLon = 179.999999
	tileMinLon = -179.999999
	tileMaxLat = 85.0511
	tileMinLat = -85.0511
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

//TileOfMerc gives the tile coordinates for the given point for the current zoom level
func (z *ZoomLevel) TileOfMerc(m PointM) Tile {
	x := math.Floor((m.E + (equator / 2)) / z.hLength)
	y := math.Floor((meridian - m.N) / z.vLength)
	t := Tile{X: int(x), Y: int(y), Z: int(z.zoom)}
	return t
}

//TileOfGeo gives the tile coordinates for the given point for the current zoom level
func (z *ZoomLevel) TileOfGeo(g PointG) (Tile, error) {
	if g.Lon < tileMinLon || g.Lon > tileMaxLon || g.Lat < tileMinLat || g.Lat > tileMaxLat {
		return Tile{}, fmt.Errorf("Point out of tiling limits: (lat, lon)(%f, %f)", g.Lat, g.Lon)
	}
	m := GeoToMerc(g)
	t := z.TileOfMerc(m)
	return t, nil
}

//ExtentOfTile return the Mercator extent of the given tile coords
func (z *ZoomLevel) ExtentOfTile(x, y int) ExtentM {
	minEast := (float64(x) * z.hLength) - (equator / 2)
	maxEast := (float64(x+1) * z.hLength) - (equator / 2)
	maxNorth := meridian - (float64(y) * z.vLength)
	minNorth := meridian - (float64(y+1) * z.vLength)
	ul := PointM{N: maxNorth, E: minEast}
	lr := PointM{N: minNorth, E: maxEast}
	me := NewExtentM(ul, lr)
	return me
}

//ExtentOf return the mercator extent of the given tile (could be slower than ZoomLevel.ExtentOfTile)
func ExtentOf(t Tile) ExtentM {
	z := NewZoomLevel(t.Z)
	return z.ExtentOfTile(t.X, t.Y)
}