package tiling_test

import (
	"math"
	"testing"

	"github.com/trealtamira/tiling"
)

func TestGeoToMerc(t *testing.T) {
	geos := []tiling.PointG{
		tiling.PointG{Lon: 9.1682557, Lat: 45.4498397},
		tiling.PointG{Lon: 2.1627174, Lat: 41.3992378},
		tiling.PointG{Lon: -123.1177155, Lat: 49.2813121},
		tiling.PointG{Lon: -68.9386812, Lat: -22.315348},
		tiling.PointG{Lon: 45.4498397, Lat: 9.1682557},
		tiling.PointG{Lon: 115.8741812, Lat: -31.8973283},
		tiling.PointG{Lon: 25.7686782, Lat: 71.1673627},
	}

	mercs := []tiling.PointM{
		tiling.PointM{N: 5692619.74164476, E: 1020605.55598653},
		tiling.PointM{N: 5071408.69855237, E: 240752.599697753},
		tiling.PointM{N: 6322729.67434633, E: -13705401.3970911},
		tiling.PointM{N: -2549428.87733578, E: -7674218.88714382},
		tiling.PointM{N: 1024989.11390234, E: 5059453.01203991},
		tiling.PointM{N: -3749840.90374462, E: 12899054.8472715},
		tiling.PointM{N: 11459741.6377346, E: 2868556.13563973},
	}
	for i := 0; i < len(geos); i++ {
		m := tiling.ToPointM(geos[i])
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
	geos := []tiling.PointG{
		tiling.PointG{Lon: 9.1682557, Lat: 45.4498397},
		tiling.PointG{Lon: 2.1627174, Lat: 41.3992378},
		tiling.PointG{Lon: -123.1177155, Lat: 49.2813121},
		tiling.PointG{Lon: -68.9386812, Lat: -22.315348},
		tiling.PointG{Lon: 45.4498397, Lat: 9.1682557},
		tiling.PointG{Lon: 115.8741812, Lat: -31.8973283},
		tiling.PointG{Lon: 25.7686782, Lat: 71.1673627},
	}

	mercs := []tiling.PointM{
		tiling.PointM{N: 5692619.74164476, E: 1020605.55598653},
		tiling.PointM{N: 5071408.69855237, E: 240752.599697753},
		tiling.PointM{N: 6322729.67434633, E: -13705401.3970911},
		tiling.PointM{N: -2549428.87733578, E: -7674218.88714382},
		tiling.PointM{N: 1024989.11390234, E: 5059453.01203991},
		tiling.PointM{N: -3749840.90374462, E: 12899054.8472715},
		tiling.PointM{N: 11459741.6377346, E: 2868556.13563973},
	}
	for i := 0; i < len(geos); i++ {
		g := tiling.ToPointG(mercs[i])
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
	mercs := []tiling.PointM{
		tiling.PointM{E: 910763.1357121654, N: 5309377.085697312},
		tiling.PointM{E: -6514065.628545966, N: -259688.542848654},
		tiling.PointM{E: 0.342789244, N: 0.342789244},
		tiling.PointM{E: 12665509.838740565, N: -3025789.0757535584},
	}
	zooms := []int{6, 4, 0, 17}

	expected := []tiling.Tile{
		tiling.Tile{X: 33, Y: 23, Z: 6},
		tiling.Tile{X: 5, Y: 8, Z: 4},
		tiling.Tile{X: 0, Y: 0, Z: 0},
		tiling.Tile{X: 106960, Y: 75432, Z: 17},
	}
	for i := 0; i < len(mercs); i++ {
		zl := tiling.NewZoomLevel(zooms[i])
		tl := zl.TileFor(mercs[i])
		te := expected[i]
		if tl.X != te.X || tl.Y != te.Y || tl.Z != te.Z {
			t.Errorf("Got the wrong tile for point %v : %v instead of %v", mercs[i], tl, te)
		}
	}
}

func TestExtent(t *testing.T) {
	tiles := []tiling.Tile{
		tiling.Tile{X: 33, Y: 23, Z: 6},
		tiling.Tile{X: 5, Y: 8, Z: 4},
		tiling.Tile{X: 0, Y: 0, Z: 0},
		tiling.Tile{X: 106960, Y: 75432, Z: 17},
	}
	exts := []tiling.MercatorExtent{
		tiling.MercatorExtent{
			SW: tiling.PointM{E: 626172.1357121654, N: 5009377.085697312},
			NE: tiling.PointM{E: 1252344.271424327, N: 5635549.221409474},
		},
		tiling.MercatorExtent{
			SW: tiling.PointM{E: -7514065.628545966, N: -2504688.542848654},
			NE: tiling.PointM{E: -5009377.085697312, N: 0},
		},
		tiling.MercatorExtent{
			SW: tiling.PointM{E: -20037508.342789244, N: -20037508.342789244},
			NE: tiling.PointM{E: 20037508.342789244, N: 20037508.342789244},
		},
		tiling.MercatorExtent{
			SW: tiling.PointM{E: 12665309.838740565, N: -3025989.0757535584},
			NE: tiling.PointM{E: 12665615.586853705, N: -3025683.327640418},
		},
	}
	for i := 0; i < len(tiles); i++ {
		ext := tiling.ExtentOf(tiles[i])
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
