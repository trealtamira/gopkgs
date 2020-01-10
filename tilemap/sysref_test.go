package tilemap_test

import (
	"math"
	"testing"

	"github.com/trealtamira/tilemap"
)

func TestGeoToMerc(t *testing.T) {
	geos := []tilemap.GeoPoint{
		tilemap.GeoPoint{Lon: 9.1682557, Lat: 45.4498397},
		tilemap.GeoPoint{Lon: 2.1627174, Lat: 41.3992378},
		tilemap.GeoPoint{Lon: -123.1177155, Lat: 49.2813121},
		tilemap.GeoPoint{Lon: -68.9386812, Lat: -22.315348},
		tilemap.GeoPoint{Lon: 45.4498397, Lat: 9.1682557},
		tilemap.GeoPoint{Lon: 115.8741812, Lat: -31.8973283},
		tilemap.GeoPoint{Lon: 25.7686782, Lat: 71.1673627},
	}

	mercs := []tilemap.MercatorPoint{
		tilemap.MercatorPoint{N: 5692619.74164476, E: 1020605.55598653},
		tilemap.MercatorPoint{N: 5071408.69855237, E: 240752.599697753},
		tilemap.MercatorPoint{N: 6322729.67434633, E: -13705401.3970911},
		tilemap.MercatorPoint{N: -2549428.87733578, E: -7674218.88714382},
		tilemap.MercatorPoint{N: 1024989.11390234, E: 5059453.01203991},
		tilemap.MercatorPoint{N: -3749840.90374462, E: 12899054.8472715},
		tilemap.MercatorPoint{N: 11459741.6377346, E: 2868556.13563973},
	}
	for i := 0; i < len(geos); i++ {
		m := tilemap.ToMercatorPoint(geos[i])
		difN := math.Abs(m.N - mercs[i].N)
		difE := math.Abs(m.E - mercs[i].E)
		if difN > 0.0000001 {
			t.Errorf("North is different (expected, actual) %f != %f %g", mercs[i].N, m.N, difN)
		}
		if difE > 0.0000001 {
			t.Errorf("East is different (expected, actual) %f != %f %f", mercs[i].E, m.E, difE)
		}
	}
}

func TestMercToGeo(t *testing.T) {
	geos := []tilemap.GeoPoint{
		tilemap.GeoPoint{Lon: 9.1682557, Lat: 45.4498397},
		tilemap.GeoPoint{Lon: 2.1627174, Lat: 41.3992378},
		tilemap.GeoPoint{Lon: -123.1177155, Lat: 49.2813121},
		tilemap.GeoPoint{Lon: -68.9386812, Lat: -22.315348},
		tilemap.GeoPoint{Lon: 45.4498397, Lat: 9.1682557},
		tilemap.GeoPoint{Lon: 115.8741812, Lat: -31.8973283},
		tilemap.GeoPoint{Lon: 25.7686782, Lat: 71.1673627},
	}

	mercs := []tilemap.MercatorPoint{
		tilemap.MercatorPoint{N: 5692619.74164476, E: 1020605.55598653},
		tilemap.MercatorPoint{N: 5071408.69855237, E: 240752.599697753},
		tilemap.MercatorPoint{N: 6322729.67434633, E: -13705401.3970911},
		tilemap.MercatorPoint{N: -2549428.87733578, E: -7674218.88714382},
		tilemap.MercatorPoint{N: 1024989.11390234, E: 5059453.01203991},
		tilemap.MercatorPoint{N: -3749840.90374462, E: 12899054.8472715},
		tilemap.MercatorPoint{N: 11459741.6377346, E: 2868556.13563973},
	}
	for i := 0; i < len(geos); i++ {
		g := tilemap.ToGeoPoint(mercs[i])
		difLa := math.Abs(g.Lat - geos[i].Lat)
		difLo := math.Abs(g.Lon - geos[i].Lon)
		if difLa > 0.0000001 {
			t.Errorf("Lat is different (expected, actual) %f != %f %g", geos[i].Lat, g.Lat, difLa)
		}
		if difLo > 0.0000001 {
			t.Errorf("Lon is different (expected, actual) %f != %f %f", geos[i].Lon, g.Lon, difLo)
		}
	}
}

func TestTileFor(t *testing.T) {
	mercs := []tilemap.MercatorPoint{
		tilemap.MercatorPoint{E: 910763.1357121654, N: 5309377.085697312},
		tilemap.MercatorPoint{E: -6514065.628545966, N: -259688.542848654},
		tilemap.MercatorPoint{E: 0.342789244, N: 0.342789244},
		tilemap.MercatorPoint{E: 12665509.838740565, N: -3025789.0757535584},
	}
	zooms := []int{6, 4, 0, 17}

	expected := []tilemap.Tile{
		tilemap.Tile{X: 33, Y: 23, Z: 6},
		tilemap.Tile{X: 5, Y: 8, Z: 4},
		tilemap.Tile{X: 0, Y: 0, Z: 0},
		tilemap.Tile{X: 106960, Y: 75432, Z: 17},
	}
	for i := 0; i < len(mercs); i++ {
		zl := tilemap.NewZoomLevel(zooms[i])
		tl := zl.TileFor(mercs[i])
		te := expected[i]
		if tl.X != te.X || tl.Y != te.Y || tl.Z != te.Z {
			t.Errorf("Got the wrong tile for point %v : %v instead of %v", mercs[i], tl, te)
		}
	}
}

func TestExtent(t *testing.T) {
	tiles := []tilemap.Tile{
		tilemap.Tile{X: 33, Y: 23, Z: 6},
		tilemap.Tile{X: 5, Y: 8, Z: 4},
		tilemap.Tile{X: 0, Y: 0, Z: 0},
		tilemap.Tile{X: 106960, Y: 75432, Z: 17},
	}
	exts := []tilemap.MercatorExtent{
		tilemap.MercatorExtent{
			SW: tilemap.MercatorPoint{E: 626172.1357121654, N: 5009377.085697312},
			NE: tilemap.MercatorPoint{E: 1252344.271424327, N: 5635549.221409474},
		},
		tilemap.MercatorExtent{
			SW: tilemap.MercatorPoint{E: -7514065.628545966, N: -2504688.542848654},
			NE: tilemap.MercatorPoint{E: -5009377.085697312, N: 0},
		},
		tilemap.MercatorExtent{
			SW: tilemap.MercatorPoint{E: -20037508.342789244, N: -20037508.342789244},
			NE: tilemap.MercatorPoint{E: 20037508.342789244, N: 20037508.342789244},
		},
		tilemap.MercatorExtent{
			SW: tilemap.MercatorPoint{E: 12665309.838740565, N: -3025989.0757535584},
			NE: tilemap.MercatorPoint{E: 12665615.586853705, N: -3025683.327640418},
		},
	}
	for i := 0; i < len(tiles); i++ {
		ext := tilemap.ExtentOf(tiles[i])
		difS := math.Abs(ext.SW.N - exts[i].SW.N)
		difN := math.Abs(ext.NE.N - exts[i].NE.N)
		difE := math.Abs(ext.NE.E - exts[i].NE.E)
		difW := math.Abs(ext.SW.E - exts[i].SW.E)
		const lim = 0.0000001
		if difS > lim || difN > lim || difW > lim || difE > lim {
			t.Errorf("Extent is different (expected, actual) %v != %v", ext, exts[i])
		}
	}

}
