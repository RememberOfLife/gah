package main

import (
	"github.com/RememberOfLife/gah"
	"github.com/fogleman/gg"
)

func main() {
	noise := gah.NewCoherentNoise(0, 0.01, 5, 2, 0.4, nil)
	const width, height int = 256, 256
	dc := gg.NewContext(width, height)

	const hueDistance float64 = 100

	for ix := 0; ix < width; ix++ {
		for iy := 0; iy < height; iy++ {
			r := noise.Eval3(float64(ix), float64(iy), 0)
			g := noise.Eval3(float64(ix), float64(iy), hueDistance)
			b := noise.Eval3(float64(ix), float64(iy), hueDistance*2)
			cvR, cvG, cvB := gah.ScaleF2I(r, -1, 1, 0, 255), gah.ScaleF2I(g, -1, 1, 0, 255), gah.ScaleF2I(b, -1, 1, 0, 255)
			dc.SetRGB255(cvR, cvG, cvB)
			dc.SetPixel(ix, iy)
		}
	}

	dc.SavePNG("./out.png")
}
