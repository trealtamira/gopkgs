package tiling

import (
	"math"
)

const (
	wgs84SphericalAxis = 6378137.0
	equator            = wgs84SphericalAxis * math.Pi * 2
	meridian           = wgs84SphericalAxis * math.Pi
	rad2deg            = 180 / math.Pi
	deg2rad            = math.Pi / 180
)

//PointM is a point in projected WebMercator coordinates (EPSG:3857) https://en.wikipedia.org/wiki/Web_Mercator_projection#Identifiers
type PointM struct {
	N float64
	E float64
}

//PointG is a point in geographic WGS84 coordinates (EPSG:4326)
type PointG struct {
	Lat float64
	Lon float64
}

//ExtentM represent a squared extent in mercator coordinates
type ExtentM struct {
	North float64
	South float64
	East  float64
	West  float64
}

// NewExtentM build a merc extent from the given opposite vertex: UL upper left, LR lower right
func NewExtentM(UL, LR PointM) ExtentM {
	ex := ExtentM{North: UL.N, South: LR.N, East: LR.E, West: UL.E}
	return ex
}

func (e ExtentM) UL() PointM {
	return PointM{N: e.North, E: e.West}
}

func (e ExtentM) LR() PointM {
	return PointM{N: e.South, E: e.East}
}

//ExtentG represent a squared extent in geographic coordinates
type ExtentG struct {
	MinLat float64
	MinLon float64
	MaxLat float64
	MaxLon float64
}

// NewExtentG build a geo extent from the given opposite vertex: UL upper left, LR lower right
func NewExtentG(UL, LR PointG) ExtentG {
	ex := ExtentG{MaxLat: UL.Lat, MinLat: LR.Lat, MaxLon: LR.Lon, MinLon: UL.Lon}
	return ex
}

func (e ExtentG) UL() PointG {
	return PointG{Lat: e.MaxLat, Lon: e.MinLon}
}

func (e ExtentG) LR() PointG {
	return PointG{Lat: e.MinLat, Lon: e.MaxLon}
}

//GeoToMerc convert the given geo point to mercator
func GeoToMerc(g PointG) PointM {
	radLat := g.Lat * deg2rad
	radLon := g.Lon * deg2rad
	north := 0.5 * math.Log((1+math.Sin(radLat))/(1-math.Sin(radLat))) * wgs84SphericalAxis
	east := radLon * wgs84SphericalAxis
	p := PointM{N: north, E: east}
	return p
}

//MercToGeo convert the given mercator point to geo
func MercToGeo(m PointM) PointG {
	lon := (m.E / wgs84SphericalAxis) * rad2deg
	lat := (2*math.Atan(math.Exp(m.N/wgs84SphericalAxis)) - 0.5*math.Pi) * rad2deg
	p := PointG{Lat: lat, Lon: lon}
	return p
}

//ToGeoExtent convert the given mercator extent to geo
func MercToGeoExt(me ExtentM) ExtentG {
	geoUL := MercToGeo(me.UL())
	geoLR := MercToGeo(me.LR())
	geoEx := NewExtentG(geoUL, geoLR)
	return geoEx
}

func Intersection(ext1, ext2 ExtentM) ExtentM {
	intersection := Extent{}
	intersection.MinX = math.Max(ext1.MinX, ext2.MinX)
	intersection.MaxX = math.Min(ext1.MaxX, ext2.MaxX)
	intersection.MinY = math.Max(ext1.MinY, ext2.MinY)
	intersection.MaxY = math.Min(ext1.MaxY, ext2.MaxY)
	ok, replacement := isConsistent(intersection)
	if !ok {
		return replacement
	}
	return intersection
}

func isConsistent(ext Extent) (bool, Extent) {
	if ext.MinX >= ext.MaxX || ext.MinY >= ext.MaxY {
		return false, Extent{MinX: 0, MinY: 0, MaxX: 0, MaxY: 0}
	}
	return true, ext

}
