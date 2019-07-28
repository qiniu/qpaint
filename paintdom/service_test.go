package paintdom

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestRoute(t *testing.T) {
	req := newRequest("POST", "http://localhost/drawings", nil)
	route, _ := getRoute(req)
	if route != "POST/drawings" {
		t.Error("TestRoute: POST /drawings")
	}

	req = newRequest("GET", "http://localhost/drawings/<DrawingID>", nil)
	route, _ = getRoute(req)
	if route != "GET/drawings/*" {
		t.Error("TestRoute: GET /drawings/<DrawingID>")
	}

	req = newRequest("GET", "http://localhost/drawings/<DrawingID>/shapes/<ShapeID>", nil)
	route, _ = getRoute(req)
	if route != "GET/drawings/*/shapes/*" {
		t.Error("TestRoute: GET /drawings/<DrawingID>/shapes/<ShapeID>")
	}
}

func newRequest(method, url string, data interface{}) *http.Request {
	b, _ := json.Marshal(data)
	body := bytes.NewReader(b)
	req, _ := http.NewRequest(method, url, body)
	return req
}
