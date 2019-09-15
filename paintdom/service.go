package paintdom

import (
	"syscall"
	"strconv"
	"strings"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"github.com/qiniu/http/httputil"
	"github.com/qiniu/http/restrpc"
	"github.com/qiniu/x/jsonutil"
)

// ---------------------------------------------------

var (
	errNotImpl  = httputil.NewError(499, "not impl")
	errBadToken = httputil.NewError(401, "bad token")
)

func mgoError(err error) error {
	if err == nil {
		return nil
	}
	if err == mgo.ErrNotFound {
		return syscall.ENOENT
	}
	if mgo.IsDup(err) {
		return syscall.EEXIST
	}
	return err
}

// ---------------------------------------------------

// Env 代表 RPC 请求的环境。
type Env struct {
	restrpc.Env
	UID UserID
}

// OpenEnv 初始化环境。
func (p *Env) OpenEnv(rcvr interface{}, w *http.ResponseWriter, req *http.Request) error {
	auth := req.Header.Get("Authorization")
	pos := strings.Index(auth, " ")
	if pos < 0 || auth[:pos] != "QPaintStub" {
		return errBadToken
	}
	uid, err := strconv.Atoi(auth[pos+1:])
	if err != nil {
		return errBadToken
	}
	p.UID = UserID(uid)
	return p.Env.OpenEnv(rcvr, w, req)
}

// ---------------------------------------------------

// Service 提供 RESTful API 层访问接口。
type Service struct {
	doc *Document
}

// NewService 创建Service实例。
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

// PostDrawingSync 同步客户端的修改。
func (p *Service) PostDrawingSync(ds *serviceDrawingSync, env *Env) (err error) {
	log.Println(env.Req.Method, env.Req.URL, jsonutil.Stringify(ds))

	changes := make([]Shape, len(ds.Changes))
	for i, item := range ds.Changes {
		changes[i] = item.Get()
	}

	id := env.Args[0]
	drawing, err := p.doc.Get(env.UID, id)
	if err != nil {
		return
	}
	return drawing.Sync(ds.Shapes, changes)
}

// PostDrawings 创建新drawing。
func (p *Service) PostDrawings(env *Env) (ret M, err error) {
	log.Println(env.Req.Method, env.Req.URL)
	drawing, err := p.doc.Add(env.UID)
	if err != nil {
		return
	}
	return M{"id": drawing.GetID()}, nil
}

// GetDrawing 取得drawing的内容。
func (p *Service) GetDrawing(env *Env) (ret M, err error) {
	log.Println(env.Req.Method, env.Req.URL)
	id := env.Args[0]
	drawing, err := p.doc.Get(env.UID, id)
	if err != nil {
		return
	}
	shapes, err := drawing.List()
	if err != nil {
		return
	}
	return M{"shapes": shapes}, nil
}

// DeleteDrawing 删除drawing。
func (p *Service) DeleteDrawing(env *Env) (err error) {
	id := env.Args[0]
	return p.doc.Delete(env.UID, id)
}

// PostShapes 创建新shape。
func (p *Service) PostShapes(aShape *serviceShape, env *Env) (err error) {
	id := env.Args[0]
	drawing, err := p.doc.Get(env.UID, id)
	if err != nil {
		return
	}
	return drawing.Add(aShape.Get())
}

// GetShape 取得一个shape的内容。
func (p *Service) GetShape(env *Env) (shape Shape, err error) {
	id := env.Args[0]
	drawing, err := p.doc.Get(env.UID, id)
	if err != nil {
		return
	}

	shapeID := env.Args[1]
	return drawing.Get(shapeID)
}

// PostShape 修改一个shape。
func (p *Service) PostShape(shapeOrZorder *serviceShapeOrZorder, env *Env) (err error) {
	id := env.Args[0]
	drawing, err := p.doc.Get(env.UID, id)
	if err != nil {
		return
	}

	shapeID := env.Args[1]
	if shapeOrZorder.Zorder != "" {
		return drawing.SetZorder(shapeID, shapeOrZorder.Zorder)
	}
	return drawing.Set(shapeID, shapeOrZorder.Get())
}

// DeleteShape 删除一个shape。
func (p *Service) DeleteShape(env *Env) (err error) {
	id := env.Args[0]
	drawing, err := p.doc.Get(env.UID, id)
	if err != nil {
		return
	}

	shapeID := env.Args[1]
	return drawing.Delete(shapeID)
}

// ---------------------------------------------------

type serviceShape struct {
	ID      string       `json:"id"`
	Path    *PathData    `json:"path,omitempty"`
	Line    *LineData    `json:"line,omitempty"`
	Rect    *RectData    `json:"rect,omitempty"`
	Ellipse *EllipseData `json:"ellipse,omitempty"`
}

func (p *serviceShape) Get() Shape {
	if p.Path != nil {
		return &Path{ShapeBase: ShapeBase{p.ID}, PathData: *p.Path}
	}
	if p.Line != nil {
		return &Line{ShapeBase: ShapeBase{p.ID}, LineData: *p.Line}
	}
	if p.Rect != nil {
		return &Rect{ShapeBase: ShapeBase{p.ID}, RectData: *p.Rect}
	}
	if p.Ellipse != nil {
		return &Ellipse{ShapeBase: ShapeBase{p.ID}, EllipseData: *p.Ellipse}
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

// Main 是 paintdom 程序的 main 入口。
func Main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	doc := NewDocument(session)
	service := NewService(doc)
	router := restrpc.Router{}
	http.ListenAndServe(":9999", router.Register(service, routeTable))
}

// ---------------------------------------------------
