package paintdom

// ---------------------------------------------------

type coord   = float64
type ShapeID = string

type Shape interface {
	GetID() ShapeID
}

type shapeBase struct {
	ID ShapeID `json:"id"`
}

func (p *shapeBase) GetID() ShapeID {
	return p.ID
}

// ---------------------------------------------------

type ShapeStyle struct {
	LineWidth coord  `json:"lineWidth"`
	LineColor string `json:"lineColor"`
	FillColor string `json:"fillColor"`
}

type Point struct {
	X coord `json:"x"`
	Y coord `json:"y"`
}

// ---------------------------------------------------

type Path struct {
	shapeBase `json:",inline"`
	pathData  `json:"path"`
}

type pathData struct {
	Points []Point    `json:"points,omitempty"`
	Close  bool       `json:"close,omitempty"`
	Style  ShapeStyle `json:"style"`
}

// ---------------------------------------------------

type Line struct {
	shapeBase `json:",inline"`
	lineData  `json:"line"`
}

type lineData struct {
	Pt1    Point      `json:"pt1"`
	Pt2    Point      `json:"pt2"`
	Style  ShapeStyle `json:"style"`
}

// ---------------------------------------------------

type Rect struct {
	shapeBase `json:",inline"`
	rectData  `json:"rect"`
}

type rectData struct {
	X      coord      `json:"x"`
	Y      coord      `json:"y"`
	Width  coord      `json:"width"`
	Height coord      `json:"height"`
	Style  ShapeStyle `json:"style"`
}

// ---------------------------------------------------

type Ellipse struct {
	shapeBase   `json:",inline"`
	ellipseData `json:"ellipse"`
}

type ellipseData struct {
	X       coord      `json:"x"`
	Y       coord      `json:"y"`
	RadiusX coord      `json:"radiusX"`
	RadiusY coord      `json:"radiusY"`
	Style   ShapeStyle `json:"style"`
}

// ---------------------------------------------------
