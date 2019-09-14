package paintdom

// ---------------------------------------------------

type coord   = float64
type ShapeID = string

type Shape interface {
	GetID() ShapeID
}

type ShapeBase struct {
	ID ShapeID `json:"id" bson:"-"`
}

func (p *ShapeBase) GetID() ShapeID {
	return p.ID
}

// ---------------------------------------------------

type ShapeStyle struct {
	LineWidth coord  `json:"lineWidth" bson:"lineWidth"`
	LineColor string `json:"lineColor" bson:"lineColor"`
	FillColor string `json:"fillColor" bson:"fillColor"`
}

type Point struct {
	X coord `json:"x" bson:"x"`
	Y coord `json:"y" bson:"y"`
}

// ---------------------------------------------------

type Path struct {
	ShapeBase `json:",inline" bson:",inline"`
	PathData  `json:"path" bson:"path"`
}

type PathData struct {
	Points []Point    `json:"points,omitempty" bson:"points,omitempty"`
	Close  bool       `json:"close,omitempty" bson:"close,omitempty"`
	Style  ShapeStyle `json:"style" bson:"style"`
}

// ---------------------------------------------------

type Line struct {
	ShapeBase `json:",inline" bson:",inline"`
	LineData  `json:"line" bson:"line"`
}

type LineData struct {
	Pt1    Point      `json:"pt1" bson:"pt1"`
	Pt2    Point      `json:"pt2" bson:"pt2"`
	Style  ShapeStyle `json:"style" bson:"style"`
}

// ---------------------------------------------------

type Rect struct {
	ShapeBase `json:",inline" bson:",inline"`
	RectData  `json:"rect" bson:"rect"`
}

type RectData struct {
	X      coord      `json:"x" bson:"x"`
	Y      coord      `json:"y" bson:"y"`
	Width  coord      `json:"width" bson:"width"`
	Height coord      `json:"height" bson:"height"`
	Style  ShapeStyle `json:"style" bson:"style"`
}

// ---------------------------------------------------

type Ellipse struct {
	ShapeBase   `json:",inline" bson:",inline"`
	EllipseData `json:"ellipse" bson:"ellipse"`
}

type EllipseData struct {
	X       coord      `json:"x" bson:"x"`
	Y       coord      `json:"y" bson:"y"`
	RadiusX coord      `json:"radiusX" bson:"radiusX"`
	RadiusY coord      `json:"radiusY" bson:"radiusY"`
	Style   ShapeStyle `json:"style" bson:"style"`
}

// ---------------------------------------------------
