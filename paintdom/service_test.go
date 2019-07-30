package paintdom

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func newServer() *httptest.Server {
	doc := NewDocument()
	service := NewService(doc)
	return httptest.NewServer(service)
}

func Post(ret interface{}, url, body string) (err error) {
	b := strings.NewReader(body)
	resp, err := http.Post(url, "application/json", b)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if (ret != nil) {
		err = json.NewDecoder(resp.Body).Decode(ret)
	}
	return
}

type idRet struct {
	ID string `json:"id"`
}

func TestNewDrawing(t *testing.T) {

	ts := newServer()
	defer ts.Close()

	var ret idRet
	err := Post(&ret, ts.URL + "/drawings", "")
	if err != nil {
		t.Fatal("Post /drawings failed:", err)
	}
	if ret.ID != "10001" {
		t.Log("new drawing id:", ret.ID)
	}
}
