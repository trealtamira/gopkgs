package tilemap

import (
	"math"
)

const (
	wgs84SphericalAxis = 6378137.0
	equator            = wgs84SphericalAxis * math.Pi * 2
	meridian           = wgs84SphericalAxis * math.Pi
	rad2deg            = 180 / math.Pi
	deg2rad            = math.Pi / 180
	maxLon             = 179.999999
	minLon             = -179.999999
	maxLat             = 85.0511
	minLat             = -85.0511
)

//MercatorPoint is a point in projected WebMercator coordinates (EPSG:3857 EPSG:900913)
type MercatorPoint struct {
	N float64
	E float64
}

//GeoPoint is a point in geographic coordinates (EPSG:4326)
type GeoPoint struct {
	Lat float64
	Lon float64
}

//Tile reresent a tile in ZoomLevel
type Tile struct {
	X int
	Y int
	Z int
}

//MercatorExtent represent a squared extent in mercator coordinates
type MercatorExtent struct {
	NE MercatorPoint
	SW MercatorPoint
}

//GeoExtent represent a squared extent in geographic coordinates
type GeoExtent struct {
	NE GeoPoint
	SW GeoPoint
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
func (z *ZoomLevel) TileFor(m MercatorPoint) Tile {
	x := math.Floor((m.E + (equator / 2)) / z.hLength)
	y := math.Floor((meridian - m.N) / z.vLength)
	t := Tile{X: int(x), Y: int(y), Z: int(z.zoom)}
	return t
}

//Extent return the extent mercator coordinates of the given tile
func (z *ZoomLevel) Extent(x, y int) MercatorExtent {
	minEast := (float64(x) * z.hLength) - (equator / 2)
	maxEast := (float64(x+1) * z.hLength) - (equator / 2)
	maxNorth := meridian - (float64(y) * z.vLength)
	minNorth := meridian - (float64(y+1) * z.vLength)
	ne := MercatorPoint{N: maxNorth, E: maxEast}
	sw := MercatorPoint{N: minNorth, E: minEast}
	me := MercatorExtent{NE: ne, SW: sw}
	return me
}

//ExtentOf return the extent mercator coordinates of the given tile
func ExtentOf(t Tile) MercatorExtent {
	z := NewZoomLevel(t.Z)
	return z.Extent(t.X, t.Y)
}

//ToMercatorPoint convert the given geo point to mercator
func ToMercatorPoint(g GeoPoint) MercatorPoint {
	radLat := g.Lat * deg2rad
	radLon := g.Lon * deg2rad
	north := 0.5 * math.Log((1+math.Sin(radLat))/(1-math.Sin(radLat))) * wgs84SphericalAxis
	east := radLon * wgs84SphericalAxis
	p := MercatorPoint{N: north, E: east}
	return p
}

//ToGeoPoint convert the given mercator point to geo
func ToGeoPoint(m MercatorPoint) GeoPoint {
	lon := (m.E / wgs84SphericalAxis) * rad2deg
	lat := (2*math.Atan(math.Exp(m.N/wgs84SphericalAxis)) - 0.5*math.Pi) * rad2deg
	p := GeoPoint{Lat: lat, Lon: lon}
	return p
}

//ToGeoExtent convert the given mercator extent to geo
func ToGeoExtent(me MercatorExtent) GeoExtent {
	geoNE := ToGeoPoint(me.NE)
	geoSW := ToGeoPoint(me.SW)
	geoEx := GeoExtent{NE: geoNE, SW: geoSW}
	return geoEx
}
