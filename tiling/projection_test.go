package tiling_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/trealtamira/gopkgs/tiling"
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
	for i, p := range geos {
		t.Run(fmt.Sprintf("Point %d (Lat,Lon)(%f, %f)", i, p.Lat, p.Lon), func(t *testing.T) {
			m := tiling.GeoToMerc(p)
			difN := math.Abs(m.N - mercs[i].N)
			difE := math.Abs(m.E - mercs[i].E)
			if difN > 0.0000001 {
				t.Errorf("North is different (expected, actual) %f != %f %g", mercs[i].N, m.N, difN)
			}
			if difE > 0.0000001 {
				t.Errorf("East is different (expected, actual) %f != %f %f", mercs[i].E, m.E, difE)
			}
		})
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
	for i, p := range mercs {
		t.Run(fmt.Sprintf("Point %d (E,N)(%f, %f)", i, p.E, p.N), func(t *testing.T) {
			g := tiling.MercToGeo(p)
			difLa := math.Abs(g.Lat - geos[i].Lat)
			difLo := math.Abs(g.Lon - geos[i].Lon)
			if difLa > 0.0000001 {
				t.Errorf("Lat is different (expected, actual) %f != %f %g", geos[i].Lat, g.Lat, difLa)
			}
			if difLo > 0.0000001 {
				t.Errorf("Lon is different (expected, actual) %f != %f %f", geos[i].Lon, g.Lon, difLo)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	exts := make(map[string]tiling.ExtentM)
	exts["ext1"] = tiling.ExtentM{West: 0, East: 60, South: 0, North: 60}
	exts["ext2"] = tiling.ExtentM{West: 60, East: 120, South: 0, North: 60}
	exts["ext3"] = tiling.ExtentM{West: 0, East: 60, South: -60, North: 0}
	exts["ext4"] = tiling.ExtentM{West: 60, East: 120, South: -60, North: 0}
	exts["extA"] = tiling.ExtentM{West: 70, East: 110, South: 20, North: 40}
	exts["extB"] = tiling.ExtentM{West: 20, East: 40, South: -20, North: 20}
	exts["extC"] = tiling.ExtentM{West: 50, East: 90, South: -50, North: -30}
	exts["in1B"] = tiling.ExtentM{West: 20, East: 40, South: 0, North: 20}
	exts["in3B"] = tiling.ExtentM{West: 20, East: 40, South: -20, North: 0}
	exts["in2A"] = tiling.ExtentM{West: 70, East: 110, South: 20, North: 40}
	exts["in3C"] = tiling.ExtentM{West: 50, East: 60, South: -50, North: -30}
	exts["in4C"] = tiling.ExtentM{West: 60, East: 90, South: -50, North: -30}
	tests := [][]string{
		{"ext1", "extB", "in1B"},
		{"ext3", "extB", "in3B"},
		{"ext2", "extA", "in2A"},
		{"ext3", "extC", "in3C"},
		{"ext4", "extC", "in4C"},
	}
	for _, e := range tests {
		t.Run(fmt.Sprintf("Intersecting %s with %s", e[0], e[1]), func(t *testing.T) {
			xs, _ := tiling.Intersection(exts[e[0]], exts[e[1]])
			xr, _ := tiling.Intersection(exts[e[1]], exts[e[0]])
			if !tiling.Equals(xs, exts[e[2]]) {
				t.Errorf("Intersection is not the expected")
			}
			if !tiling.Equals(xr, exts[e[2]]) {
				t.Errorf("Intersection reverse is not the expected")
			}
		})
	}
}

func TestIntersectionEmpty(t *testing.T) {
	exts := make(map[string]tiling.ExtentM)
	exts["ext1"] = tiling.ExtentM{West: 0, East: 60, South: 0, North: 60}
	exts["ext2"] = tiling.ExtentM{West: 60, East: 120, South: 0, North: 60}
	exts["ext3"] = tiling.ExtentM{West: 0, East: 60, South: -60, North: 0}
	exts["ext4"] = tiling.ExtentM{West: 60, East: 120, South: -60, North: 0}
	exts["extA"] = tiling.ExtentM{West: 70, East: 110, South: 20, North: 40}
	exts["extB"] = tiling.ExtentM{West: 20, East: 40, South: -20, North: 20}
	exts["extC"] = tiling.ExtentM{West: 50, East: 90, South: -50, North: -30}
	tests := [][]string{
		{"ext1", "extA"},
		{"ext3", "extA"},
		{"ext1", "extC"},
		{"ext4", "extB"},
		{"ext2", "extB"},
		{"ext2", "extC"},
	}
	for _, e := range tests {
		t.Run(fmt.Sprintf("Intersecting %s with %s", e[0], e[1]), func(t *testing.T) {
			xs, oks := tiling.Intersection(exts[e[0]], exts[e[1]])
			xr, okr := tiling.Intersection(exts[e[1]], exts[e[0]])
			if oks {
				t.Errorf("Intersection %+v should be empty", xs)
			}
			if okr {
				t.Errorf("Intersection %+v should be empty", xr)
			}
		})
	}
}
