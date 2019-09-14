package paintdom

import (
	"testing"

	"github.com/qiniu/http/restrpc"
	"github.com/qiniu/qiniutest/httptest"
	"github.com/qiniu/x/mockhttp"
	"gopkg.in/mgo.v2"
)

func newTestingDocument(t *testing.T) *Document {
	DBName = "testQPaint"
	session, err := mgo.Dial("localhost")
	if err != nil {
		t.Fatal("mgo.Dial failed:", err)
	}
	session.DB(DBName).DropDatabase()
	return NewDocument(session)
}

func TestService(t *testing.T) {

	doc := newTestingDocument(t)
	service := NewService(doc)

	transport := mockhttp.NewTransport()
	router := restrpc.Router{}
	transport.ListenAndServe("qpaint.com", router.Register(service, routeTable))

	ctx := httptest.New(t)
	ctx.SetTransport(transport)

	ctx.Exec(
	`
	post http://qpaint.com/drawings
	header Authorization 'QPaintStub 1'
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

	match $(line1v2) '{
		"id": "1",
		"line": {
			"pt1": {"x":152,"y":333},
			"pt2": {"x":158,"y":324},
			"style": {"lineWidth":1,"lineColor":"black"}
		}
	}'

	match $(line2) '{
		"id": "2",
		"line": {
			"pt1": {"x":152,"y":133},
			"pt2": {"x":358,"y":324},
			"style": {"lineWidth":1,"lineColor":"black","fillColor":"white"}
		}
	}'

	post http://qpaint.com/drawings/$(id1)/shapes
	header Authorization 'QPaintStub 1'
	json $(line1)
	ret 200

	get http://qpaint.com/drawings/$(id1)/shapes/1
	header Authorization 'QPaintStub 1'
	ret 200
	json $(line1)

	post http://qpaint.com/drawings/$(id1)/sync
	header Authorization 'QPaintStub 1'
	json '{
		"changes": [$(line1v2), $(line2)],
		"shapes": ["1", "2"]
	}'
	ret 200

	get http://qpaint.com/drawings/$(id1)/shapes/1
	header Authorization 'QPaintStub 1'
	ret 200
	json $(line1v2)

	get http://qpaint.com/drawings/$(id1)/shapes/2
	header Authorization 'QPaintStub 1'
	ret 200
	json $(line2)
	`)
}
