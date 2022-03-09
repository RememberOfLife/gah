package gah

import (
	"github.com/fogleman/gg"
)

//TODO calculate proper encoding of checksum and use correct start stop bar
// https://de.qwe.wiki/wiki/RM4SCC
// https://www.morovia.com/kb/Royal-Mail-Barcode-RMS4CC-10630.html

//TODO make a pixel based version of this

// DrawSeedBarCode1 draws a SBC1 in the RM4SCC style at the desired location calculated as following:
//
// anchor point is `x - w * ax`, `y - h * ay`
//
// use `ax=0.5`, `ay=0.5` to center to a point and both at zero so the code is drawn to the bottom right of the anchor
//
// for best resolution use w in multiples of 112 and h in multiples of 3
func DrawSeedBarCode1(dc *gg.Context, data uint64, x float64, y float64, w float64, h float64, ax float64, ay float64) {
	dc.Push()

	var strokeWidth float64 = w / (38 + 37*2) // room for 38bars and 37spaces at double the width

	dc.Translate(x-w*ax, y-h*ay)

	dc.SetLineWidth(strokeWidth)
	// start and end bar are fulheight missing the step-mark
	dc.DrawLine(0, 0, 0, h*1/3)
	dc.Stroke()
	dc.DrawLine(0, h*2/3, 0, h)
	dc.Stroke()
	dc.Translate(strokeWidth*3, 0)

	checkTop, checkBottom := 0, 0
	for ix := 35; ix >= 0; ix-- { // 35 iterations because uint64 is processed in 32 pairs of 2 bits each plus one checksum byte
		var lineType uint64
		if ix > 3 {
			// processing the uint64
			lineType = (^(data >> ((ix - 4) * 2))) & 0b11 // move the relevant 2 bits to the relevant positions and single them out
			if lineType&0b10 == 0 {
				checkTop++
			}
			if lineType&0b01 == 0 {
				checkBottom++
			}
		} else {
			// processing the fake checksum byte
			lineType = 0b11
			if (checkTop>>ix)&0b1 == 1 {
				lineType &= 0b01
			}
			if (checkBottom>>ix)&0b1 == 1 {
				lineType &= 0b10
			}
		}
		switch lineType {
		case 0:
			dc.DrawLine(0, 0, 0, h) // 0=00
		case 1:
			dc.DrawLine(0, 0, 0, h*2/3) // 1==01
		case 2:
			dc.DrawLine(0, h*1/3, 0, h) // 2=10
		case 3:
			dc.DrawLine(0, h*1/3, 0, h*2/3) // 3=11
		}
		dc.Stroke()
		dc.Translate(strokeWidth*3, 0)
	}
	// start and end bar are fulheight missing the step-mark
	dc.DrawLine(0, 0, 0, h*1/3)
	dc.Stroke()
	dc.DrawLine(0, h*2/3, 0, h)
	dc.Stroke()

	dc.Pop()
}
