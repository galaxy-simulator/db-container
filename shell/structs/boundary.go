package structs

type Boundary struct {
	// boundary box values
	x     int64
	y     int64
	width int64
}

func NewBoundary(x int64, y int64, width int64) *Boundary {
	return &Boundary{x: x, y: y, width: width}
}

func (b *Boundary) Width() int64 {
	return b.width
}

func (b *Boundary) SetWidth(width int64) {
	b.width = width
}

func (b *Boundary) Y() int64 {
	return b.y
}

func (b *Boundary) SetY(y int64) {
	b.y = y
}

func (b *Boundary) X() int64 {
	return b.x
}

func (b *Boundary) SetX(x int64) {
	b.x = x
}
