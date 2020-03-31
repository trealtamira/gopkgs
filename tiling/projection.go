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
	maxLon             = 179.999999
	minLon             = -179.999999
	maxLat             = 85.0511
	minLat             = -85.0511
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
	NE PointM
	SW PointM
}

//ExtentG represent a squared extent in geographic coordinates
type ExtentG struct {
	NE PointG
	SW PointG
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
	geoNE := MercToGeo(me.NE)
	geoSW := MercToGeo(me.SW)
	geoEx := ExtentG{NE: geoNE, SW: geoSW}
	return geoEx
}
