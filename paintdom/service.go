package paintdom

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"io"
	"strconv"
	"strings"
	"syscall"
)

// ---------------------------------------------------

type M map[string]interface{}
type RouteTable map[string]func(w http.ResponseWriter, req *http.Request, args []string)

type Service struct {
	doc        *Document
	routeTable RouteTable
}

func NewService(doc *Document) (p *Service) {
	p = &Service{doc: doc}
	p.routeTable = RouteTable{
		"POST/drawings":              p.PostDrawings,
		"GET/drawings/*":             p.GetDrawing,
		"DELETE/drawings/*":          p.DeleteDrawing,
		"POST/drawings/*/sync":       p.PostDrawingSync,
		"POST/drawings/*/shapes":     p.PostShapes,
		"GET/drawings/*/shapes/*":    p.GetShape,
		"POST/drawings/*/shapes/*":   p.PostShape,
		"DELETE/drawings/*/shapes/*": p.DeleteShape,
	}
	return
}

func (p *Service) PostDrawingSync(w http.ResponseWriter, req *http.Request, args []string) {
	b := bytes.NewBuffer(nil)
	io.Copy(b, req.Body)
	log.Println(req.Method, req.URL, b.String())

	var ds serviceDrawingSync
	err := json.NewDecoder(b).Decode(&ds)
	if err != nil {
		ReplyError(w, err)
		return
	}

	changes := make([]Shape, len(ds.Changes))
	for i, item := range ds.Changes {
		changes[i] = item.Get()
	}

	id := args[0]
	drawing, err := p.doc.Get(id)
	if err != nil {
		ReplyError(w, err)
		return
	}
	err = drawing.Sync(ds.Shapes, changes)
	if err != nil {
		ReplyError(w, err)
		return
	}
	ReplyCode(w, 200)
}

func (p *Service) PostDrawings(w http.ResponseWriter, req *http.Request, args []string) {
	log.Println(req.Method, req.URL)
	drawing, err := p.doc.Add()
	if err != nil {
		ReplyError(w, err)
		return
	}
	Reply(w, 200, M{"id": drawing.ID})
}

func (p *Service) GetDrawing(w http.ResponseWriter, req *http.Request, args []string) {
	log.Println(req.Method, req.URL)
	id := args[0]
	drawing, err := p.doc.Get(id)
	if err != nil {
		ReplyError(w, err)
		return
	}
	shapes, err := drawing.List()
	if err != nil {
		ReplyError(w, err)
		return
	}
	Reply(w, 200, M{"shapes": shapes})
}

func (p *Service) DeleteDrawing(w http.ResponseWriter, req *http.Request, args []string) {
	id := args[0]
	err := p.doc.Delete(id)
	if err != nil {
		ReplyError(w, err)
		return
	}
	ReplyCode(w, 200)
}

func (p *Service) PostShapes(w http.ResponseWriter, req *http.Request, args []string) {
	id := args[0]
	drawing, err := p.doc.Get(id)
	if err != nil {
		ReplyError(w, err)
		return
	}

	var aShape serviceShape
	err = json.NewDecoder(req.Body).Decode(&aShape)
	if err != nil {
		ReplyError(w, err)
		return
	}

	err = drawing.Add(aShape.Get())
	if err != nil {
		ReplyError(w, err)
		return
	}
	ReplyCode(w, 200)
}

func (p *Service) GetShape(w http.ResponseWriter, req *http.Request, args []string) {
	id := args[0]
	drawing, err := p.doc.Get(id)
	if err != nil {
		ReplyError(w, err)
		return
	}

	shapeID := args[1]
	shape, err := drawing.Get(shapeID)
	if err != nil {
		ReplyError(w, err)
		return
	}
	Reply(w, 200, shape)
}

func (p *Service) PostShape(w http.ResponseWriter, req *http.Request, args []string) {
	id := args[0]
	drawing, err := p.doc.Get(id)
	if err != nil {
		ReplyError(w, err)
		return
	}

	var shapeID = args[1]
	var shapeOrZorder serviceShapeOrZorder
	err = json.NewDecoder(req.Body).Decode(&shapeOrZorder)
	if err != nil {
		ReplyError(w, err)
		return
	}

	if shapeOrZorder.Zorder != "" {
		err = drawing.SetZorder(shapeID, shapeOrZorder.Zorder)
	} else {
		err = drawing.Set(shapeID, shapeOrZorder.Get())
	}
	if err != nil {
		ReplyError(w, err)
		return
	}
	ReplyCode(w, 200)
}

func (p *Service) DeleteShape(w http.ResponseWriter, req *http.Request, args []string) {
	id := args[0]
	drawing, err := p.doc.Get(id)
	if err != nil {
		ReplyError(w, err)
		return
	}

	shapeID := args[1]
	err = drawing.Delete(shapeID)
	if err != nil {
		ReplyError(w, err)
		return
	}
	ReplyCode(w, 200)
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

func Reply(w http.ResponseWriter, code int, data interface{}) {
	b, _ := json.Marshal(data)
	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(code)
	w.Write(b)
	log.Println("REPLY", code, string(b))
}

func ReplyCode(w http.ResponseWriter, code int) {
	header := w.Header()
	header.Set("Content-Length", "0")
	w.WriteHeader(code)
	log.Println("REPLY", code)
}

func ReplyError(w http.ResponseWriter, err error) {
	if err == syscall.ENOENT {
		Reply(w, 404, M{"error": "entry not found"})
	} else if err == syscall.EINVAL {
		Reply(w, 400, M{"error": "invalid arguments"})
	} else if err == syscall.EEXIST {
		Reply(w, 409, M{"error": "entry already exists"})
	} else {
		Reply(w, 500, M{"error": err.Error()})
	}
}

// ---------------------------------------------------

func (p *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route, args := getRoute(req)
	if handle, ok := p.routeTable[route]; ok {
		handle(w, req, args)
	}
}

func getRoute(req *http.Request) (route string, args []string) {
	parts := strings.Split(req.URL.Path, "/")
	parts[0] = req.Method
	for i := 2; i < len(parts); i += 2 {
		args = append(args, parts[i])
		parts[i] = "*"
	}
	route = strings.Join(parts, "/")
	return
}

// ---------------------------------------------------

func Main() {
	doc := NewDocument()
	service := NewService(doc)
	http.ListenAndServe(":9999", service)
}

// ---------------------------------------------------
