package main

import (
	"image/color"

	"github.com/RememberOfLife/gah"
	"github.com/fogleman/gg"
)

func main() {
	var wPx int = 1000
	var hPx int = 1000
	dc := gg.NewContext(wPx, hPx)

	gradient := gah.ColorRamp{[]gah.ColorStop{
		gah.ColorStop{0, color.RGBA{0, 0, 0, 0xFF}},
		gah.ColorStop{0.25, color.RGBA{255, 0, 0, 0xFF}},
		gah.ColorStop{0.5, color.RGBA{0, 255, 0, 0xFF}},
		gah.ColorStop{0.75, color.RGBA{0, 0, 255, 0xFF}},
		gah.ColorStop{1, color.RGBA{255, 255, 255, 0xFF}},
	}}

	for ix := 0; ix < wPx; ix++ {
		for iy := 0; iy < hPx; iy++ {
			dc.SetColor(gradient.Sample(float64(ix) / float64(wPx)))
			dc.SetPixel(ix, iy)
		}
	}

	dc.SavePNG("./out.png")
}
