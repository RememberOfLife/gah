package gah

import (
	"image/color"
	"sort"
)

// ColorRamp holds multiple colorstops between which can be interpolated
// use Sort to order them by their position
type ColorRamp struct {
	GradientStops []ColorStop // must be sorted
}

// ColorStop defines the position at which a color is strongest in the ColorRamp
type ColorStop struct {
	Position float64 // in range [0, 1]
	Color    color.RGBA
}

// Sample returns the value on the ColorRamp gradient that is calculated at the given position
// interpolates linearly between the given color stops
// using this on a ColorRamp with unsorted stops may break
func (cr *ColorRamp) Sample(position float64) color.RGBA {
	if len(cr.GradientStops) < 2 {
		return color.RGBA{0, 0, 0, 0xFF}
	}
	if position < 0 || position < cr.GradientStops[0].Position {
		return cr.GradientStops[0].Color
	}
	if position > 1 || position > cr.GradientStops[len(cr.GradientStops)-1].Position {
		return cr.GradientStops[len(cr.GradientStops)-1].Color
	}
	// find the 2 relevant stop for this position
	gradientIndex := 0
	var cs1 ColorStop = cr.GradientStops[gradientIndex]
	var cs2 ColorStop = cr.GradientStops[gradientIndex+1]
	for position > cs2.Position {
		gradientIndex++
		cs1 = cr.GradientStops[gradientIndex]
		cs2 = cr.GradientStops[gradientIndex+1]
	}
	// interpolate using RGBMix
	return RGBMix(cs1.Color, cs2.Color, ScaleF2F(position, cs1.Position, cs2.Position, 0, 1), true)
}

// Sort sorts the ColorStops of a ColorRamp by their position so Sample does not break
func (cr *ColorRamp) Sort() {
	sort.Slice(cr.GradientStops, func(i, j int) bool {
		return cr.GradientStops[i].Position < cr.GradientStops[j].Position
	})
}
