package gah

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

type quadTreeQuadrant int

const (
	topRight    quadTreeQuadrant = 0
	topLeft     quadTreeQuadrant = 1
	bottomLeft  quadTreeQuadrant = 2
	bottomRight quadTreeQuadrant = 3
)

// QuadTreeLeafThreshold defines the maximum amount of points allowed in a QuadTree leafnode before it is split
const QuadTreeLeafThreshold = 10

// QuadTree is a typical 2D QuadTree containing points
type QuadTree struct {
	x, y, w, h     float64 // position and extent of the quad
	subTrees       [4]*QuadTree
	leafPoints     [QuadTreeLeafThreshold]Vec2f
	leafPointCount int // counts the number of leafs points, internal nodes accumulate their childrens counts
}

// NewQuadTree creates a new QuadTree with the given dimensions
func NewQuadTree(x float64, y float64, w float64, h float64) (qt *QuadTree) {
	return &QuadTree{x, y, w, h, [4]*QuadTree{}, [QuadTreeLeafThreshold]Vec2f{}, 0}
}

func (qt *QuadTree) isLeaf() bool {
	for _, st := range qt.subTrees {
		if st != nil {
			return false
		}
	}
	return true
}

// InsertPoint inserts the given point into the QuadTree, branching the tree where necessary
//TODO make this iterative
func (qt *QuadTree) InsertPoint(p Vec2f) {
	if ok, _ := qt.Contains(p); !ok {
		return // do nothing if point if outside of tree
	}
	// point is inside of tree
	if qt.isLeaf() {
		// is leaf
		if qt.leafPointCount < QuadTreeLeafThreshold {
			// insert into leaf with space
			qt.leafPoints[qt.leafPointCount] = p
			qt.leafPointCount++
			return
		}
		// set node to internal, and insert points into generated subtrees
		//TODO sparsely generate subtrees
		hw := qt.w / 2
		hh := qt.h / 2
		qt.subTrees[topRight] = NewQuadTree(qt.x+hw, qt.y, hw, hh)
		qt.subTrees[topLeft] = NewQuadTree(qt.x, qt.y, hw, hh)
		qt.subTrees[bottomLeft] = NewQuadTree(qt.x, qt.y+hh, hw, hh)
		qt.subTrees[bottomRight] = NewQuadTree(qt.x+hw, qt.y+hh, hw, hh)
		for _, lp := range qt.leafPoints {
			_, quad := qt.Contains(lp)
			qt.subTrees[quad].InsertPoint(lp)
		}
		_, quad := qt.Contains(p)
		qt.subTrees[quad].InsertPoint(p)
		return
	}
	// is internal, iterate tree until bounding leaf is found, insert there
	cqt := qt
	for !cqt.isLeaf() {
		_, quad := qt.Contains(p)
		cqt = cqt.subTrees[quad]
	}
	cqt.InsertPoint(p)
}

// InsertPoints inserts the given points into the QuadTree, branching the tree where necessary
func (qt *QuadTree) InsertPoints(p []Vec2f) {
	for _, np := range p {
		qt.InsertPoint(np)
	}
}

func (qt *QuadTree) Contains(p Vec2f) (bool, quadTreeQuadrant) {
	if p.X < qt.x && p.X > qt.x+qt.w && p.Y < qt.y && p.Y > qt.y+qt.h {
		return false, -1
	}
	xsign := int(math.Copysign(1, (qt.x+qt.w/2)-p.X)+1) / 2
	ysign := int(math.Copysign(1, (qt.y+qt.h/2)-p.Y)+1) / 2
	ysign = ysign ^ 1 // graphics coordinate system is other way around, so switch y axis
	return true, quadTreeQuadrant((xsign ^ ysign) + ysign + ysign)
}

func (qt *QuadTree) Intersects(x float64, y float64, w float64, h float64) bool {
	return !(qt.x+qt.w < x || x+w < qt.x || qt.y+qt.h < y || y+h < qt.y)
}

//TODO make this iterative
func (qt *QuadTree) GetPoints() (results []Vec2f) {
	if qt.isLeaf() {
		return qt.leafPoints[:qt.leafPointCount]
	}
	for _, st := range qt.subTrees {
		results = append(results, st.GetPoints()...)
	}
	return nil
}

// SignedDistanceToPoint returns a relative distance value from p to the QuadTrees bounding box
// returns >0 if outside the box, <0 if inside, and 0 when exactly on the line
func (qt *QuadTree) SignedDistanceToPoint(p Vec2f) float64 {
	/* https://stackoverflow.com/questions/30545052/calculate-signed-distance-between-point-and-rectangle#30545544
	float sdAxisAlignedRect(vec2 uv, vec2 tl, vec2 br)
	{
		vec2 d = max(tl-uv, uv-br);
		return length(max(vec2(0.0), d)) + min(0.0, max(d.x, d.y));
	}*/
	d := Vec2f{
		math.Max(qt.x-p.X, p.X-(qt.x+qt.w)),
		math.Max(qt.y-p.Y, p.Y-(qt.y+qt.h)),
	}
	l := math.Hypot(
		math.Max(0, d.X),
		math.Max(0, d.Y),
	)
	return l + math.Min(0, math.Max(d.X, d.Y))
}

// QueryRange returns all leaf points of QuadTree inside the given region
//TODO make this iterative
func (qt *QuadTree) QueryRange(x float64, y float64, w float64, h float64) (results []Vec2f) {
	//TODO optimization: if a subtree is entirely contained in the range, skip checking and append all its combined points
	results = []Vec2f{}
	if !qt.Intersects(x, y, w, h) {
		// return empty list if range doesnt intersect this tree
		return
	}
	if qt.isLeaf() {
		// is leaf, check and add children where neccessary
		for _, p := range qt.leafPoints[:qt.leafPointCount] {
			if !(p.X < qt.x && p.X > qt.x+qt.w && p.Y < qt.y && p.Y > qt.y+qt.h) {
				// point contained, append
				results = append(results, p)
			}
		}
		return
	}
	for _, st := range qt.subTrees {
		results = append(results, st.QueryRange(x, y, w, h)...)
	}
	return
}

// QueryKNN returns the k nearest neighbors to the given point p
//TODO make this iterative
func (qt *QuadTree) QueryKNN(p Vec2f, k int) []Vec2f {
	/*//TODO use a proper priority queue
	type prioTree struct {
		priority float64
		tree     *QuadTree
	}
	ptqPush := func(ptq *[]prioTree, npt prioTree) {
		for i, pt := range *ptq {
			if pt.priority > npt.priority {
				*ptq = append((*ptq)[:i+1], (*ptq)[i:]...)
				(*ptq)[i] = npt
				return
			}
		}
		*ptq = append(*ptq, npt)
	}
	ptqPop := func(ptq *[]prioTree) (pt prioTree) {
		pt = (*ptq)[0]
		*ptq = (*ptq)[1:]
		return
	}*/
	return nil //TODO
}

//DEBUG REMOVE
func DrawQuadTree(dc *gg.Context, tree *QuadTree) {
	// draw border os my quad
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(1)
	dc.DrawRectangle(tree.x, tree.y, tree.w, tree.h)
	dc.Stroke()
	if tree.isLeaf() {
		for i, p := range tree.leafPoints {
			if i >= tree.leafPointCount {
				break
			}
			i++
			// draw point
			dc.SetRGB(1, 0, 0)
			dc.DrawCircle(p.X, p.Y, 2)
			dc.Fill()
		}
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(fmt.Sprintf("%d", tree.leafPointCount), tree.x+tree.w/2, tree.y+tree.h/2, 0.5, 0.5)
		return
	}
	for _, st := range tree.subTrees {
		DrawQuadTree(dc, st)
	}
}
