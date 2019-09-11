package paintdom

import (
	"testing"

	"github.com/qiniu/http/restrpc"
	"github.com/qiniu/qiniutest/httptest"
	"github.com/qiniu/x/mockhttp"
)

func TestNewDrawing(t *testing.T) {

	doc := NewDocument()
	service := NewService(doc)

	transport := mockhttp.NewTransport()
	router := restrpc.Router{}
	transport.ListenAndServe("qpaint.com", router.Register(service, routeTable))

	ctx := httptest.New(t)
	ctx.SetTransport(transport)

	ctx.Exec(
	`
	post http://qpaint.com/drawings
	ret 200
	json '{
		"id": $(id1)
	}'

	match $(line1) '{
		"id": "1",
		"line": {
			"pt1": {"x": 2.0, "y": 3.0},
			"pt2": {"x": 15.0, "y": 30.0},
			"style": {
				"lineWidth": 3,
				"lineColor": "red"
			}
		}
	}'

	post http://qpaint.com/drawings/$(id1)/shapes
	json $(line1)
	ret 200

	get http://qpaint.com/drawings/$(id1)/shapes/1
	ret 200
	json $(line1)
	`)

	if !ctx.GetVar("id1").Equal("10001") {
		t.Fatal(`$(id1) != "10001"`)
	}
}
