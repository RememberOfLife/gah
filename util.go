package gah

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

// Vec2i is a simple 2D int vector
type Vec2i struct {
	X, Y int
}

// Vec2f is a simple 2D float vector
type Vec2f struct {
	X, Y float64
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ImgFastSaveToPNG skips png compression for much faster saving speed at the cost of more memory
func ImgFastSaveToPNG(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := &png.Encoder{
		CompressionLevel: png.NoCompression,
	}
	return enc.Encode(file, img)
}

// ImageToRGBA converts an image to an RGBA image, hopefully using a lower level go construct for better performance
// ripped straight from fogleman/gg
func ImageToRGBA(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

// ImgGet returns a [0, 1] grayscale value representing the color on the given image at the given coordinates
func ImgGet(img image.Image, x, y int) float64 {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	if x < 0 || x >= w || y < 0 || y >= h {
		return 0
	}
	r, _, _, _ := color.GrayModel.Convert(img.At(x, y)).RGBA()
	return float64(r & 0xFF)
}

// ImgGetRGBA returns the rgba components of the color on the given image at the given coordinates
func ImgGetRGBA(img image.Image, x, y int) (r, g, b, a int) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	if x < 0 || x >= w || y < 0 || y >= h {
		return 0, 0, 0, 0
	}
	tr, tg, tb, ta := img.At(x, y).RGBA()
	return int(tr & 0xFF), int(tg & 0xFF), int(tb & 0xFF), int(ta & 0xFF)
}

// RGBMix linearly interpolates between the two given colors, alpha may be maxed to 0xFF to prevent decay
func RGBMix(color1 color.RGBA, color2 color.RGBA, ratio2 float64, maxAlpha bool) color.RGBA {
	r1 := float64(color1.R) / 255
	g1 := float64(color1.G) / 255
	b1 := float64(color1.B) / 255
	a1 := float64(color1.A) / 255
	r2 := float64(color2.R) / 255
	g2 := float64(color2.G) / 255
	b2 := float64(color2.B) / 255
	a2 := float64(color2.A) / 255
	var a uint8
	if maxAlpha {
		a = 0xFF
	} else {
		a = uint8(((1-ratio2)*a1 + ratio2*a2) * 255)
	}
	return color.RGBA{
		uint8(((1-ratio2)*r1 + ratio2*r2) * 255),
		uint8(((1-ratio2)*g1 + ratio2*g2) * 255),
		uint8(((1-ratio2)*b1 + ratio2*b2) * 255),
		a,
	}
}

// MixF mixes the provided values together using the given ratio for val2
func MixF(val1 float64, val2 float64, ratio2 float64) float64 {
	return (1-ratio2)*val1 + ratio2*val2
}

// MixI works as MixF does but on integers
func MixI(val1 int, val2 int, ratio2 float64) int {
	return int(MixF(float64(val1), float64(val2), ratio2))
}

// ScaleF2F returns a number in the interval [outMin, outMax] that is percentage-wise as much removed from each end of the interval as inNum in [inMin, inMax]
func ScaleF2F(inNum float64, inMin float64, inMax float64, outMin float64, outMax float64) float64 {
	return (inNum-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

// ScaleF2I returns a number in the interval [outMin, outMax] that is percentage-wise as much removed from each end of the interval as inNum in [inMin, inMax]
func ScaleF2I(inNum float64, inMin float64, inMax float64, outMin int, outMax int) int {
	return int((inNum-inMin)*(float64(outMax-outMin))/(inMax-inMin) + float64(outMin))
}

// ScaleI2F returns a number in the interval [outMin, outMax] that is percentage-wise as much removed from each end of the interval as inNum in [inMin, inMax]
func ScaleI2F(inNum int, inMin int, inMax int, outMin float64, outMax float64) float64 {
	return float64(inNum-inMin)*(outMax-outMin)/float64(inMax-inMin) + outMin
}

// ScaleI2I returns a number in the interval [outMin, outMax] that is percentage-wise as much removed from each end of the interval as inNum in [inMin, inMax]
func ScaleI2I(inNum int, inMin int, inMax int, outMin int, outMax int) int {
	return int((inNum-inMin)*(outMax-outMin)/(inMax-inMin) + outMin)
}

// Clamp returns a number inside [inMin, inMax]
func Clamp(inNum float64, inMin float64, inMax float64) float64 {
	if inNum > inMax {
		return inMax
	}
	if inNum < inMin {
		return inMin
	}
	return inNum
}
