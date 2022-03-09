package gah

import (
	"math"
	"math/rand"
	"sort"

	"github.com/fogleman/poissondisc"
)

// VoronoiDiagram2D represents a 2D voronoi diagram
type VoronoiDiagram2D struct {
	Seed       uint64
	X, Y, W, H float64
	Points     []Vec2f
	Scale      float64
	K          int
	PdsTrys    int
}

// NewVoronoiDiagram2D creates a new voronoi diagram, with points spaced to have a minimum distance given by the scale
// if k is -1 then crackle will be used instead of k nearest neighbor, i.e. return distance to nearest edge
func NewVoronoiDiagram2D(seed uint64, x float64, y float64, w float64, h float64, scale float64, k int, pdsTrys int) *VoronoiDiagram2D {
	vd := &VoronoiDiagram2D{seed, x, y, w, h, []Vec2f{}, scale, k, pdsTrys}
	for _, sample := range poissondisc.Sample(x-scale, y-scale, w+scale, h+scale, scale, pdsTrys, rand.New(rand.NewSource(int64(seed)))) {
		vd.Points = append(vd.Points, Vec2f{sample.X, sample.Y})
	}
	return vd
}

// GetParamSignature returns a byte slice containing all relevant unique parameters
func (vd *VoronoiDiagram2D) GetParamSignature() (signature []byte) {
	signature = append(signature, IntToBytes(int(vd.Seed))...)
	signature = append(signature, Float64ToBytes(vd.X)...)
	signature = append(signature, Float64ToBytes(vd.Y)...)
	signature = append(signature, Float64ToBytes(vd.W)...)
	signature = append(signature, Float64ToBytes(vd.H)...)
	signature = append(signature, Float64ToBytes(vd.Scale)...)
	signature = append(signature, IntToBytes(vd.K)...)
	signature = append(signature, IntToBytes(vd.PdsTrys)...)
	return signature
}

// GetEvalRange returns the min and max values that can be expected from the Eval2
func (vd *VoronoiDiagram2D) GetEvalRange() (outMin float64, outMax float64) {
	return 0, 1
}

// Eval2 returns the distance to the k nearest neighbor
// returns within range [0, 1]; or 0 for out of bounds; 1 is closest to a point
func (vd *VoronoiDiagram2D) Eval2(x, y float64) float64 {
	if x < vd.X || x >= vd.X+vd.W || y < vd.Y || y >= vd.Y+vd.H {
		return 0
	}
	// calc dist to every point
	type distIndex struct {
		i    int
		dist float64
	}
	var distances []distIndex
	for i, p := range vd.Points {
		distances = append(distances, distIndex{i, math.Hypot(x-p.X, y-p.Y)})
	}
	// sort distances
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].dist < distances[j].dist
	})
	k := vd.K
	var borderDist float64 = distances[k+1].dist
	var targetDist float64 = 0
	if k == -1 { // if crackle
		k = 0
		borderDist = (distances[0].dist + distances[1].dist) / 2
		targetDist = distances[1].dist - distances[0].dist
	}
	targetDist = distances[k].dist
	// return k nearest
	return 1 - ScaleF2F(targetDist, 0, borderDist+1, 0, 1)
}
