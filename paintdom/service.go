package paintdom

import (
	"log"
	"net/http"

	"github.com/qiniu/http/restrpc"
)

// ---------------------------------------------------

type M map[string]interface{}

type Service struct {
	doc *Document
}

func NewService(doc *Document) (p *Service) {
	p = &Service{doc: doc}
	return
}

var routeTable = [][2]string{
	{"POST /drawings", "PostDrawings"},
	{"GET /drawings/*", "GetDrawing"},
	{"DELETE /drawings/*", "DeleteDrawing"},
	{"POST /drawings/*/sync", "PostDrawingSync"},
	{"POST /drawings/*/shapes", "PostShapes"},
	{"GET /drawings/*/shapes/*", "GetShape"},
	{"POST /drawings/*/shapes/*", "PostShape"},
	{"DELETE /drawings/*/shapes/*", "DeleteShape"},
}

func (p *Service) PostDrawingSync(ds *serviceDrawingSync, env *restrpc.Env) (err error) {
	log.Println(env.Req.Method, env.Req.URL, *ds)

	changes := make([]Shape, len(ds.Changes))
	for i, item := range ds.Changes {
		changes[i] = item.Get()
	}

	id := env.Args[0]
	err = p.doc.Sync(id, ds.Shapes, changes)
	return
}

func (p *Service) PostDrawings(env *restrpc.Env) (ret M, err error) {
	log.Println(env.Req.Method, env.Req.URL)
	drawing, err := p.doc.Add()
	if err != nil {
		return
	}
	return M{"id": drawing.ID}, nil
}

func (p *Service) GetDrawing(env *restrpc.Env) (ret M, err error) {
	log.Println(env.Req.Method, env.Req.URL)
	id := env.Args[0]
	drawing, err := p.doc.Get(id)
	if err != nil {
		return
	}
	shapes, err := drawing.List()
	if err != nil {
		return
	}
	return M{"shapes": shapes}, nil
}

func (p *Service) DeleteDrawing(env *restrpc.Env) (err error) {
	id := env.Args[0]
	return p.doc.Delete(id)
}

func (p *Service) PostShapes(aShape *serviceShape, env *restrpc.Env) (err error) {
	id := env.Args[0]
	drawing, err := p.doc.Get(id)
	if err != nil {
		return
	}
	return drawing.Add(aShape.Get())
}

func (p *Service) GetShape(env *restrpc.Env) (shape Shape, err error) {
	id := env.Args[0]
	drawing, err := p.doc.Get(id)
	if err != nil {
		return
	}

	shapeID := env.Args[1]
	return drawing.Get(shapeID)
}

func (p *Service) PostShape(shapeOrZorder *serviceShapeOrZorder, env *restrpc.Env) (err error) {
	id := env.Args[0]
	drawing, err := p.doc.Get(id)
	if err != nil {
		return
	}

	shapeID := env.Args[1]
	if shapeOrZorder.Zorder != "" {
		return drawing.SetZorder(shapeID, shapeOrZorder.Zorder)
	}
	return drawing.Set(shapeID, shapeOrZorder.Get())
}

func (p *Service) DeleteShape(env *restrpc.Env) (err error) {
	id := env.Args[0]
	drawing, err := p.doc.Get(id)
	if err != nil {
		return
	}

	shapeID := env.Args[1]
	return drawing.Delete(shapeID)
}

// ---------------------------------------------------

type serviceShape struct {
	ID      string       `json:"id"`
	Path    *pathData    `json:"path"`
	Line    *lineData    `json:"line"`
	Rect    *rectData    `json:"rect"`
	Ellipse *ellipseData `json:"ellipse"`
}

func (p *serviceShape) Get() Shape {
	if p.Path != nil {
		return &Path{shapeBase: shapeBase{p.ID}, pathData: *p.Path}
	}
	if p.Line != nil {
		return &Line{shapeBase: shapeBase{p.ID}, lineData: *p.Line}
	}
	if p.Rect != nil {
		return &Rect{shapeBase: shapeBase{p.ID}, rectData: *p.Rect}
	}
	if p.Ellipse != nil {
		return &Ellipse{shapeBase: shapeBase{p.ID}, ellipseData: *p.Ellipse}
	}
	return nil
}

type serviceShapeOrZorder struct {
	serviceShape `json:",inline"`
	Zorder       string `json:"zorder"`
}

type serviceDrawingSync struct {
	Changes []serviceShape `json:"changes"`
	Shapes  []ShapeID      `json:"shapes"`
}

// ---------------------------------------------------

func Main() {
	doc := NewDocument()
	service := NewService(doc)
	router := restrpc.Router{}
	http.ListenAndServe(":9999", router.Register(service, routeTable))
}

// ---------------------------------------------------
