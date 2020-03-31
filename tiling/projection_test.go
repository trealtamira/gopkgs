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
		t.Run(fmt.Sprintf("Testing point %d", i), func(t *testing.T) {
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
		t.Run(fmt.Sprintf("Testing point %d", i), func(t *testing.T) {
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
