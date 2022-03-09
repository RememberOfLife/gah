package main

import (
	"math/rand"

	"github.com/RememberOfLife/gah"
	"github.com/fogleman/gg"
)

func main() {
	const width, height int = 1000, 1000
	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	qt := gah.NewQuadTree(0, 0, 1000, 1000)

	for i := 0; i < 50; i++ {
		p := gah.Vec2f{
			float64(width) * gah.Clamp(rand.NormFloat64()*0.1+0.5, 0, 1),
			float64(height) * gah.Clamp(rand.NormFloat64()*0.1+0.5, 0, 1),
		}
		qt.InsertPoint(p)
	}

	gah.DrawQuadTree(dc, qt)

	dc.SavePNG("./out.png")
}
