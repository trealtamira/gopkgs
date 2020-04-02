package tiling_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/trealtamira/gopkgs/tiling"
)

func TestZoomLevelTileOfMerc(t *testing.T) {
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
	for i, p := range mercs {
		t.Run(fmt.Sprintf("Point %d (E,N)(%f, %f)", i, p.E, p.N), func(t *testing.T) {
			zl := tiling.NewZoomLevel(zooms[i])
			tl := zl.TileOfMerc(p)
			te := expected[i]
			if tl.X != te.X || tl.Y != te.Y || tl.Z != te.Z {
				t.Errorf("Got the wrong tile for point %v : %v instead of %v", mercs[i], tl, te)
			}
		})
	}
}

func TestExtentOf(t *testing.T) {
	tiles := []tiling.Tile{
		tiling.Tile{X: 33, Y: 23, Z: 6},
		tiling.Tile{X: 5, Y: 8, Z: 4},
		tiling.Tile{X: 0, Y: 0, Z: 0},
		tiling.Tile{X: 106960, Y: 75432, Z: 17},
	}
	exts := []tiling.ExtentM{
		tiling.ExtentM{
			West: 626172.1357121654, South: 5009377.085697312, East: 1252344.271424327, North: 5635549.221409474,
		},
		tiling.ExtentM{
			West: -7514065.628545966, South: -2504688.542848654, East: -5009377.085697312, North: 0,
		},
		tiling.ExtentM{
			West: -20037508.342789244, South: -20037508.342789244, East: 20037508.342789244, North: 20037508.342789244,
		},
		tiling.ExtentM{
			West: 12665309.838740565, South: -3025989.0757535584, East: 12665615.586853705, North: -3025683.327640418,
		},
	}
	for i, e := range tiles {
		t.Run(fmt.Sprintf("Tile %d (X,Y,Z)(%d, %d, %d)", i, e.X, e.Y, e.Z), func(t *testing.T) {
			ext := tiling.ExtentOf(e)
			difS := math.Abs(ext.South - exts[i].South)
			difN := math.Abs(ext.North - exts[i].North)
			difE := math.Abs(ext.East - exts[i].East)
			difW := math.Abs(ext.West - exts[i].West)
			const lim = 0.0000001
			if difS > lim || difN > lim || difW > lim || difE > lim {
				t.Errorf("Extent is different (expected, actual) %+v != %+v", ext, exts[i])
			}
		})
	}
}
