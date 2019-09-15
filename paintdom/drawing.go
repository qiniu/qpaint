package paintdom

import (
	"syscall"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserID 是用户ID，代表一个用户。
type UserID = uint64

// M 只是一个缩略写法。
type M = bson.M

// DBName 是默认的 QPaint Database Name。
var DBName = "qpaint" 

// ---------------------------------------------------

// QShape 是 drawing 返回的 shape 实例。
type QShape map[string]interface{}

// GetID 取得该 shape 的 id。
func (shape QShape) GetID() string {
	if val, ok1 := shape["id"]; ok1 {
		if id, ok2 := val.(string); ok2 {
			return id
		}
	}
	return ""
}

func (shape QShape) setID(id string) {
	shape["id"] = id
}

// ---------------------------------------------------

// Drawing 代表用户的一个drawing文档。
type Drawing struct {
	id bson.ObjectId
	session *mgo.Session
}

func newDrawing(id bson.ObjectId, session *mgo.Session) *Drawing {
	return &Drawing{
		id: id,
		session: session,
	}
}

// GetID 取得 drawing ID。
func (p *Drawing) GetID() string {
	return p.id.Hex()
}

// Sync 同步 drawing 的修改。
func (p *Drawing) Sync(shapes []ShapeID, changes []Shape) (err error) {
	c := p.session.Copy()
	defer c.Close()
	shapeColl := c.DB(DBName).C("shape")
	for _, change := range changes {
		spid := change.GetID()
		_, err = shapeColl.Upsert(M{
			"dgid": p.id,
			"spid": spid, 
		}, M{
			"dgid": p.id,
			"spid": spid, 
			"shape": change,
		})
		if err != nil {
			return mgoError(err)
		}
	}
	drawingColl := c.DB(DBName).C("drawing")
	return mgoError(drawingColl.UpdateId(p.id, M{
		"$set": M{"shapes": shapes},
	}))
}

// Add 添加新图形。
func (p *Drawing) Add(shape Shape) (err error) {
	c := p.session.Copy()
	defer c.Close()
	shapeColl := c.DB(DBName).C("shape")
	spid := shape.GetID()
	err = shapeColl.Insert(M{
		"dgid": p.id,
		"spid": spid,
		"shape": shape,
	})
	if err != nil {
		return mgoError(err)
	}
	drawingColl := c.DB(DBName).C("drawing")
	return mgoError(drawingColl.UpdateId(p.id, M{
		"$push": M{"shapes": spid},
	}))
}

// List 列出所有图形。
func (p *Drawing) List() (shapes []Shape, err error) {
	c := p.session.Copy()
	defer c.Close()
	var result []struct{
		ID    string `bson:"spid"`
		Shape QShape `bson:"shape"`
	}
	shapeColl := c.DB(DBName).C("shape")
	err = shapeColl.Find(M{
		"dgid": p.id,
	}).Select(M{
		"spid": 1, "shape": 1,
	}).All(&result)
	if err != nil {
		return nil, mgoError(err)
	}
	shapes = make([]Shape, len(result))
	for i, item := range result {
		item.Shape.setID(item.ID)
		shapes[i] = item.Shape
	}
	return
}

// Get 取出某个图形。
func (p *Drawing) Get(id ShapeID) (shape Shape, err error) {
	c := p.session.Copy()
	defer c.Close()
	var o struct {
		Shape QShape `bson:"shape"`
	}
	shapeColl := c.DB(DBName).C("shape")
	err = shapeColl.Find(M{
		"dgid": p.id,
		"spid": id,
	}).Select(M{
		"shape": 1,
	}).One(&o)
	if err != nil {
		return nil, mgoError(err)
	}
	o.Shape.setID(id)
	return o.Shape, nil
}

// Set 修改某个图形。
func (p *Drawing) Set(id ShapeID, shape Shape) (err error) {
	if shape.GetID() != "" {
		return syscall.EINVAL
	}
	c := p.session.Copy()
	defer c.Close()
	shapeColl := c.DB(DBName).C("shape")
	return mgoError(shapeColl.Update(M{
		"dgid": p.id,
		"spid": id,
	}, M{
		"$set": M{"shape": shape},
	}))
}

// SetZorder 修改图形的图层。
func (p *Drawing) SetZorder(id ShapeID, zorder string) (err error) {
	return errNotImpl
}

// Delete 删除某个图形。
func (p *Drawing) Delete(id ShapeID) (err error) {
	c := p.session.Copy()
	defer c.Close()
	drawingColl := c.DB(DBName).C("drawing")
	err = drawingColl.UpdateId(p.id, M{
		"$pull": M{"shapes": id},
	})
	if err != nil {
		return mgoError(err)
	}
	shapeColl := c.DB(DBName).C("shape")
	return mgoError(shapeColl.Remove(M{
		"dgid": p.id,
		"spid": id,
	}))
}

// ---------------------------------------------------

// Document 代表整个 QPaint DOM 的根。
type Document struct {
	session *mgo.Session
}

// NewDocument 创建一个 QPaint DOM 对象。
func NewDocument(session *mgo.Session) *Document {
	drawingColl := session.DB(DBName).C("drawing")
	drawingColl.EnsureIndex(mgo.Index{
		Key: []string{"uid"},
	})
	shape := session.DB(DBName).C("shape")
	shape.EnsureIndex(mgo.Index{
		Key: []string{"dgid", "spid"},
		Unique: true,
		DropDups: true,
	})
	return &Document{
		session: session,
	}
}

// Add 创建新drawing。
func (p *Document) Add(uid UserID) (drawing *Drawing, err error) {
	c := p.session.Copy()
	defer c.Close()
	drawingColl := c.DB(DBName).C("drawing")
	id := bson.NewObjectId()
	err = drawingColl.Insert(M{
		"_id": id,
		"uid": uid,
		"shapes": []ShapeID{},
	})
	if err != nil {
		return
	}
	return newDrawing(id, p.session), nil
}

// Get 获取drawing。
// 我们会检查要获取的drawing是否为该uid所拥有，如果不属于则获取会失败。
func (p *Document) Get(uid UserID, dgid string) (drawing *Drawing, err error) {
	c := p.session.Copy()
	defer c.Close()
	drawingColl := c.DB(DBName).C("drawing")
	id := bson.ObjectIdHex(dgid)
	n, err := drawingColl.Find(M{
		"_id": id,
		"uid": uid,
	}).Count()
	if err != nil {
		return nil, err
	}
	if n > 0 {
		return newDrawing(id, p.session), nil
	}
	return nil, syscall.ENOENT
}

// Delete 删除drawing。
// 我们会检查要删除的drawing是否为该uid所拥有，如果不属于删除会失败。
func (p *Document) Delete(uid UserID, dgid string) (err error) {
	c := p.session.Copy()
	defer c.Close()
	id := bson.ObjectIdHex(dgid)
	drawingColl := c.DB(DBName).C("drawing")
	err = drawingColl.Remove(M{
		"_id": id,
		"uid": uid,
	})
	if err != nil {
		return mgoError(err)
	}
	shapeColl := c.DB(DBName).C("shape")
	_, err = shapeColl.RemoveAll(M{
		"dgid": id,
	})
	return mgoError(err)
}

// ---------------------------------------------------
