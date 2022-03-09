package main

import (
	"github.com/RememberOfLife/gah"
	"github.com/fogleman/gg"
)

func main() {
	var wPx int = 1000
	var hPx int = 1000
	dc := gg.NewContext(wPx, hPx)

	voronoi := gah.NewVoronoiDiagram2D(0, 0, 0, float64(wPx), float64(hPx), 300, -1, 30)

	// draw distances
	for iy := 0; iy < hPx; iy++ {
		for ix := 0; ix < wPx; ix++ {
			c := voronoi.Eval2(float64(ix), float64(iy))
			dc.SetRGB(c, c, c)
			dc.SetPixel(ix, iy)
		}
	}

	// overlay original points of voronoi
	dc.SetRGB(1, 0, 0)
	dc.SetLineWidth(1)
	for _, sample := range voronoi.Points {
		dc.DrawPoint(sample.X, sample.Y, 2)
		dc.Fill()
	}

	dc.SavePNG("./out.png")
}
