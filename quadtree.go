package quadtree

const nodeCapacity = 4

type XY struct {
	X, Y float64
}

func NewXY(x, y float64) XY {
	return XY{x, y}
}

type AABB struct {
	center, halfDim XY
}

func NewAABB(center, halfDim XY) *AABB {
	return &AABB{center, halfDim}
}

func (a *AABB) ContainsPoint(p XY) bool {
	return p.X >= a.center.X-a.halfDim.X &&
		p.X <= a.center.X+a.halfDim.X &&
		p.Y >= a.center.Y-a.halfDim.Y &&
		p.Y <= a.center.Y+a.halfDim.Y
}

func (a *AABB) IntersectsAABB(other *AABB) bool {
	return !(other.center.X+other.halfDim.X < a.center.X-a.halfDim.X ||
		other.center.Y+other.halfDim.Y < a.center.Y-a.halfDim.Y ||
		other.center.X-other.halfDim.X > a.center.X+a.halfDim.X ||
		other.center.Y-other.halfDim.Y > a.center.Y+a.halfDim.Y)
}

type QuadTree struct {
	boundary       AABB
	points         []XY
	ul, ur, dl, dr *QuadTree
}

func (qt *QuadTree) isLeaf() bool {
	return qt.ul == nil
}

func New(boundary AABB) *QuadTree {
	points := make([]XY, 0, nodeCapacity)
	qt := &QuadTree{boundary: boundary, points: points}
	return qt
}

func (qt *QuadTree) Insert(p XY) bool {
	if !qt.boundary.ContainsPoint(p) {
		return false
	}

	if len(qt.points) < cap(qt.points) {
		qt.points = append(qt.points, p)
		return true
	}

	if qt.isLeaf() {
		qt.split()
	}

	_ = qt.ul.Insert(p) ||
		qt.ur.Insert(p) ||
		qt.dl.Insert(p) ||
		qt.dr.Insert(p)
	return true
}

func (qt *QuadTree) split() {
	if !qt.isLeaf() {
		return
	}

	boxUL := AABB{
		XY{qt.boundary.center.X - qt.boundary.halfDim.X/2,
			qt.boundary.center.Y + qt.boundary.halfDim.Y/2},
		XY{qt.boundary.halfDim.X / 2, qt.boundary.halfDim.Y / 2}}
	boxUR := AABB{
		XY{qt.boundary.center.X + qt.boundary.halfDim.X/2,
			qt.boundary.center.Y + qt.boundary.halfDim.Y/2},
		XY{qt.boundary.halfDim.X / 2, qt.boundary.halfDim.Y / 2}}
	boxDL := AABB{
		XY{qt.boundary.center.X - qt.boundary.halfDim.X/2,
			qt.boundary.center.Y - qt.boundary.halfDim.Y/2},
		XY{qt.boundary.halfDim.X / 2, qt.boundary.halfDim.Y / 2}}
	boxDR := AABB{
		XY{qt.boundary.center.X + qt.boundary.halfDim.X/2,
			qt.boundary.center.Y - qt.boundary.halfDim.Y/2},
		XY{qt.boundary.halfDim.X / 2, qt.boundary.halfDim.Y / 2}}

	qt.ul = New(boxUL)
	qt.ur = New(boxUR)
	qt.dl = New(boxDL)
	qt.dr = New(boxDR)

	for _, p := range qt.points {
		_ = qt.ul.Insert(p) ||
			qt.ur.Insert(p) ||
			qt.dl.Insert(p) ||
			qt.dr.Insert(p)
	}
	qt.points = nil
}

func (qt *QuadTree) SearchArea(a *AABB) []XY {
	results := make([]XY, 0, nodeCapacity)
	if !qt.boundary.IntersectsAABB(a) {
		return results
	}

	if qt.isLeaf() {
		for _, p := range qt.points {
			if a.ContainsPoint(p) {
				results = append(results, p)
			}
		}
		return results
	}

	results = append(results, qt.ul.SearchArea(a)...)
	results = append(results, qt.ur.SearchArea(a)...)
	results = append(results, qt.dl.SearchArea(a)...)
	results = append(results, qt.dr.SearchArea(a)...)
	return results
}

func (qt *QuadTree) SearchNear(p XY, d float64) []XY {
	d2 := d * d
	box := boundingBox(p, d)
	candidates := qt.SearchArea(box)
	results := make([]XY, 0, len(candidates))
	for _, p2 := range candidates {
		if distance2(p, p2) <= d2 {
			results = append(results, p2)
		}
	}
	return results
}

func (qt *QuadTree) IsAnyPointArea(a *AABB) bool {
	if !qt.boundary.IntersectsAABB(a) {
		return false
	}

	if qt.isLeaf() {
		for _, p := range qt.points {
			if a.ContainsPoint(p) {
				return true
			}
		}
		return false
	}

	return qt.ul.IsAnyPointArea(a) ||
		qt.ur.IsAnyPointArea(a) ||
		qt.dl.IsAnyPointArea(a) ||
		qt.dr.IsAnyPointArea(a)
}

func (qt *QuadTree) IsAnyPointNear(p XY, d float64) bool {
	a := boundingBox(p, d)
	if !qt.boundary.IntersectsAABB(a) {
		return false
	}

	if qt.isLeaf() {
		d2 := d * d
		for _, p2 := range qt.points {
			if a.ContainsPoint(p2) && distance2(p, p2) <= d2 {
				return true
			}
		}
		return false
	}

	return qt.ul.IsAnyPointNear(p, d) ||
		qt.ur.IsAnyPointNear(p, d) ||
		qt.dl.IsAnyPointNear(p, d) ||
		qt.dr.IsAnyPointNear(p, d)
}

func distance2(p1, p2 XY) float64 {
	d1 := p1.X - p2.X
	d2 := p2.Y - p2.Y
	return (d1 * d1) + (d2 * d2)
}

func boundingBox(p XY, d float64) *AABB {
	return &AABB{p, XY{d, d}}
}
