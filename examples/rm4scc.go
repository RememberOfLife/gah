package main

import (
	"github.com/RememberOfLife/gah"
	"github.com/fogleman/gg"
)

func main() {
	const wr = 112
	const hr = 3
	dc := gg.NewContext(wr*3+20, hr*20+20)
	dc.SetRGB255(255, 255, 255)
	dc.Clear()

	var data uint64 = 0x0123456789ABCDEF

	dc.SetRGB(0, 0, 0)
	gah.DrawSeedBarCode1(dc, data, 10, 10, wr*3, hr*20, 0, 0)

	dc.SavePNG("./out.png")
}
