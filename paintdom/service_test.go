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
		"id": "10001"
	}'
	`)
}
