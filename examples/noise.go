package main

import (
	"github.com/RememberOfLife/gah"
	"github.com/fogleman/gg"
)

func main() {
	noise := gah.NewCoherentNoise(0, 0.005, 5, 2, 0.5, nil)
	const width, height int = 1000, 1000
	dc := gg.NewContext(width, height)

	for ix := 0; ix < width; ix++ {
		for iy := 0; iy < height; iy++ {
			r := noise.Eval2(float64(ix), float64(iy))
			cv := gah.ScaleF2I(r, -1, 1, 0, 255)
			dc.SetRGB255(cv, cv, cv)
			dc.SetPixel(ix, iy)
		}
	}

	dc.SavePNG("./out.png")
}
