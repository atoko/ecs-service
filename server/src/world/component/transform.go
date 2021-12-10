package component

type Transform struct {
	X float32
	Y float32
	Width int
	Height int
}

func (t *Transform) GetTransform() *Transform {
	return t
}

type TransformFace interface {
	GetTransform() *Transform
}

type V2 struct {
	X, Y float32
}

func (a V2) Vector2() V2 { return a }

func (a V2) Dot(v V2) float32 {
	b := v
	return a.X*b.X + a.Y*b.Y
}

func (a V2) Sub(v V2) V2 {
	b := v
	return V2{a.X - b.X, a.Y - b.Y}
}