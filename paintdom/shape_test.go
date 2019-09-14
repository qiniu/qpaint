package paintdom

import (
	"encoding/json"
	"reflect"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func bsonMarshal(val interface{}) (b2 []byte, err error) {
	b, err := bson.Marshal(val)
	if err != nil {
		return
	}
	val2 := reflect.New(reflect.TypeOf(val)).Interface()
	err = bson.Unmarshal(b, val2)
	if err != nil {
		return
	}
	return json.Marshal(val2)
}

func TestPathEncode(t *testing.T) {
	val := new(Path)
	val.Points = []Point{
		{100, 200},
		{200, 300},
	}
	b, _ := json.Marshal(val)
	if string(b) != `{"id":"","path":{"points":[{"x":100,"y":200},{"x":200,"y":300}],"style":{"lineWidth":0,"lineColor":"","fillColor":""}}}` {
		t.Error("unexpected: json.Marshal path:", string(b))
	}
}

func TestPathBsonEncode(t *testing.T) {
	val := new(Path)
	val.Points = []Point{
		{100, 200},
		{200, 300},
	}
	b, err := bsonMarshal(val)
	if string(b) != `{"id":"","path":{"points":[{"x":100,"y":200},{"x":200,"y":300}],"style":{"lineWidth":0,"lineColor":"","fillColor":""}}}` {
		t.Error("unexpected: bson.Marshal path:", string(b), "error:", err)
	}
}

func TestRectEncode(t *testing.T) {
	val := new(Rect)
	b, _ := json.Marshal(val)
	if string(b) != `{"id":"","rect":{"x":0,"y":0,"width":0,"height":0,"style":{"lineWidth":0,"lineColor":"","fillColor":""}}}` {
		t.Error("unexpected: json.Marshal rect:", string(b))
	}
}

func TestRectBsonEncode(t *testing.T) {
	val := new(Rect)
	b, _ := bsonMarshal(val)
	if string(b) != `{"id":"","rect":{"x":0,"y":0,"width":0,"height":0,"style":{"lineWidth":0,"lineColor":"","fillColor":""}}}` {
		t.Error("unexpected: bson.Marshal rect:", string(b))
	}
}

func TestEllipseEncode(t *testing.T) {
	val := new(Ellipse)
	b, _ := json.Marshal(val)
	if string(b) != `{"id":"","ellipse":{"x":0,"y":0,"radiusX":0,"radiusY":0,"style":{"lineWidth":0,"lineColor":"","fillColor":""}}}` {
		t.Error("unexpected: json.Marshal ellipse:", string(b))
	}
}

func TestEllipseBsonEncode(t *testing.T) {
	val := new(Ellipse)
	b, _ := bsonMarshal(val)
	if string(b) != `{"id":"","ellipse":{"x":0,"y":0,"radiusX":0,"radiusY":0,"style":{"lineWidth":0,"lineColor":"","fillColor":""}}}` {
		t.Error("unexpected: bson.Marshal ellipse:", string(b))
	}
}
