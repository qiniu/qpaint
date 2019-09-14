package paintdom

import (
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
)

// ---------------------------------------------------

type shapeOnDrawing struct {
	front *shapeOnDrawing
	back  *shapeOnDrawing
	data  Shape
}

func (p *shapeOnDrawing) init() {
	p.front, p.back = p, p
}

func (p *shapeOnDrawing) insertFront(shape *shapeOnDrawing) {
	shape.back = p
	shape.front = p.front
	p.front.back = shape
	p.front = shape
}

func (p *shapeOnDrawing) insertBack(shape *shapeOnDrawing) {
	shape.front = p
	shape.back = p.back
	p.back.front = shape
	p.back = shape
}

func (p *shapeOnDrawing) moveFront() {
	front := p.front
	p.delete()
	front.insertFront(p)
}

func (p *shapeOnDrawing) moveBack() {
	back := p.back
	p.delete()
	back.insertBack(p)
}

func (p *shapeOnDrawing) delete() {
	p.front.back = p.back
	p.back.front = p.front
	p.back, p.front = nil, nil
}

// ---------------------------------------------------

type Drawing struct {
	ID     string
	idBase int64

	mutex  sync.Mutex
	shapes map[ShapeID]*shapeOnDrawing
	list   shapeOnDrawing
}

func newDrawing(id string) *Drawing {
	p := &Drawing{
		ID: id,
		shapes: make(map[ShapeID]*shapeOnDrawing),
	}
	p.list.init()
	return p
}

func (p *Drawing) Sync(shapes []ShapeID, changes []Shape) (err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	dgshapes := make([]*shapeOnDrawing, len(shapes))
	for i, id := range shapes {
		if dgshape, ok := p.shapes[id]; ok {
			dgshape.delete()
			dgshapes[i] = dgshape
		} else {
			dgshape := new(shapeOnDrawing)
			dgshapes[i] = dgshape
			p.shapes[id] = dgshape
		}
	}
	head := &p.list
	for item := head.front; item != head; item = item.front {
		delete(p.shapes, item.data.GetID())
	}
	head.init()
	for _, dgshape := range dgshapes {
		head.insertBack(dgshape)
	}
	for _, change := range changes {
		id := change.GetID()
		p.shapes[id].data = change
	}
	return
}

func (p *Drawing) Add(shape Shape) (err error) {
	id := shape.GetID()
	dgshape := &shapeOnDrawing{
		data: shape,
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.shapes[id]; ok {
		return syscall.EEXIST
	}
	p.list.insertBack(dgshape)
	p.shapes[id] = dgshape
	return
}

func (p *Drawing) List() (shapes []Shape, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	n := len(p.shapes)
	shapes = make([]Shape, n)
	item := p.list.front
	for i := 0; i < n; i++ {
		shapes[i] = item.data
		item = item.front
	}
	return
}

func (p *Drawing) Get(id ShapeID) (shape Shape, err error) {
	if dgshape, ok := p.shapes[id]; ok {
		return dgshape.data, nil
	}
	return nil, syscall.ENOENT
}

func (p *Drawing) Set(id ShapeID, shape Shape) (err error) {
	if shape.GetID() != id {
		return syscall.EINVAL
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if dgshape, ok := p.shapes[id]; ok {
		dgshape.data = shape
		return nil
	}
	return syscall.ENOENT
}

func (p *Drawing) SetZorder(id ShapeID, zorder string) (err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if shape, ok := p.shapes[id]; ok {
		switch zorder {
		case "top":
			shape.delete()
			p.list.insertBack(shape)
		case "bottom":
			shape.delete()
			p.list.insertFront(shape)
		case "front":
			if shape.front != &p.list {
				shape.moveFront()
			}
		case "back":
			if shape.back != &p.list {
				shape.moveBack()
			}
		default:
			return syscall.EINVAL
		}
		return nil
	}
	return syscall.ENOENT
}

func (p *Drawing) Delete(id ShapeID) (err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	dgshape, ok := p.shapes[id]
	if !ok {
		return syscall.ENOENT
	}
	dgshape.delete()
	delete(p.shapes, id)
	return
}

// ---------------------------------------------------

type Document struct {
	mutex sync.Mutex
	data  map[string]*Drawing
}

func NewDocument() *Document {
	drawings := make(map[string]*Drawing)
	return &Document{
		data: drawings,
	}
}

func (p *Document) Add() (drawing *Drawing, err error) {
	id := makeDrawingID()
	drawing = newDrawing(id)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.data[id] = drawing
	return
}

func (p *Document) Get(id string) (drawing *Drawing, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	drawing, ok := p.data[id]
	if !ok {
		return nil, syscall.ENOENT
	}
	return
}

func (p *Document) Delete(id string) (err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	delete(p.data, id)
	return
}

// ---------------------------------------------------

var (
	idDrawingBase int64 = 10000
)

func makeDrawingID() string {
	id := atomic.AddInt64(&idDrawingBase, 1)
	return strconv.Itoa(int(id))
}

// ---------------------------------------------------
