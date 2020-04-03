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

//UL Upper Left vertex
func (e ExtentM) UL() PointM {
	return PointM{N: e.North, E: e.West}
}

//UR Upper Right vertex
func (e ExtentM) UR() PointM {
	return PointM{N: e.North, E: e.East}
}

//LR Lower Right vertex
func (e ExtentM) LR() PointM {
	return PointM{N: e.South, E: e.East}
}

//LL Lower Left vertex
func (e ExtentM) LL() PointM {
	return PointM{N: e.South, E: e.West}
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

//UL Upper Left vertex
func (e ExtentG) UL() PointG {
	return PointG{Lat: e.MaxLat, Lon: e.MinLon}
}

//UR Upper Right vertex
func (e ExtentG) UR() PointG {
	return PointG{Lat: e.MaxLat, Lon: e.MaxLon}
}

//LR Lower Right vertex
func (e ExtentG) LR() PointG {
	return PointG{Lat: e.MinLat, Lon: e.MaxLon}
}

//LL Lower Left vertex
func (e ExtentG) LL() PointG {
	return PointG{Lat: e.MinLat, Lon: e.MinLon}
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

//MercToGeoExt convert the given mercator extent to geo
func MercToGeoExt(me ExtentM) ExtentG {
	geoUL := MercToGeo(me.UL())
	geoLR := MercToGeo(me.LR())
	geoEx := NewExtentG(geoUL, geoLR)
	return geoEx
}

//Intersection compute the intersection between two extent and false if it is empty
func Intersection(ext1, ext2 ExtentM) (ExtentM, bool) {
	ix := ExtentM{}
	ix.West = math.Max(ext1.West, ext2.West)
	ix.East = math.Min(ext1.East, ext2.East)
	ix.South = math.Max(ext1.South, ext2.South)
	ix.North = math.Min(ext1.North, ext2.North)
	if ix.West >= ix.East || ix.South >= ix.North {
		return ExtentM{West: 0, South: 0, East: 0, North: 0}, false
	}
	return ix, true
}

//Equals returns false if not equals
func Equals(ext1, ext2 ExtentM) bool {
	if ext1.North != ext2.North {
		return false
	} else if ext1.South != ext2.South {
		return false
	} else if ext1.East != ext2.East {
		return false
	} else if ext1.West != ext2.West {
		return false
	} else {
		return true
	}
}
