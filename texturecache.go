package gah

import (
	"crypto"
	_ "crypto/sha512" // register hash function
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/fogleman/gg"
)

// TextureCachable signifies a texture that can be precomputed and cached for later use
type TextureCachable interface {
	GetParamSignature() (signature []byte)
	GetEvalRange() (outMin float64, outMax float64)
	Eval2(x float64, y float64) float64
}

// IntToBytes returns a byte slice representing the input
func IntToBytes(i int) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(i))
	return buf[:]
}

// Float64ToBytes returns a byte slice representing the input
func Float64ToBytes(f float64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

// TextureCache represents a cached texture
type TextureCache struct {
	x, y, w, h int
	texture    *image.RGBA
}

// NewTextureCache returns a new cached texture provider, path should only be a directory specification (including the trailing '/')
// this method will block and pre-generate the whole texture
func NewTextureCache(provider TextureCachable, x int, y int, w int, h int, path string) *TextureCache {
	hasher := crypto.SHA512.New()
	hasher.Write(provider.GetParamSignature())
	// write own size params
	hasher.Write(IntToBytes(x))
	hasher.Write(IntToBytes(y))
	hasher.Write(IntToBytes(w))
	hasher.Write(IntToBytes(h))
	paramHash := MB64E.EncodeToString(hasher.Sum(nil)[:48]) // first 384 bits of  512 bit hash for 64 characters of base64
	cachePath := fmt.Sprintf("%s%s.png", path, paramHash)
	bounds := image.Rect(x, y, w, h)
	imgData := image.NewRGBA(bounds)
	if fileExists(cachePath) {
		imgRaw, _ := gg.LoadPNG(cachePath)
		draw.Draw(imgData, bounds, imgRaw, bounds.Min, draw.Src)
		return &TextureCache{x, y, w, h, imgData}
	}
	var tc *TextureCache = &TextureCache{x, y, w, h, imgData}
	emin, emax := provider.GetEvalRange()
	for ix := x; ix < x+w; ix++ {
		for iy := y; iy < y+h; iy++ {
			ixf := float64(ix)
			iyf := float64(iy)
			c := uint8(ScaleF2I(provider.Eval2(ixf, iyf), emin, emax, 0, 255))
			tc.texture.SetRGBA(ix, iy, color.RGBA{c, c, c, 0xFF})
		}
	}
	ImgFastSaveToPNG(tc.texture, cachePath)
	return tc
}

// Sample returns the grayscale value [0, 1] at the given position in the cached texture
func (tc *TextureCache) Sample(x int, y int) float64 {
	if x < tc.x || x >= tc.x+tc.w || y < tc.y || y >= tc.y+tc.h {
		return 0
	}
	return float64(tc.texture.RGBAAt(x, y).R) / 255
}
