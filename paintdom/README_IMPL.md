QPaint DOM IMPL (第 42 讲)
========

## 逻辑 DOM 结构

```
<Drawing1>, 隶属于某个<uid>
    <Shape11>
    ...
    <Shape1M>
...
<DrawingN>, 隶属于某个<uid>
```

## 使用界面（接口）

### Document 类

```
Add(uid UserID) (drawing *Drawing, err error)
```
* 创建新drawing。 

```
Get(uid UserID, id string) (drawing *Drawing, err error)
```
* 获取drawing。
* 我们会检查要获取的drawing是否为该uid所拥有，如果不属于则获取会失败。 

```
Delete(uid UserID, id string) (err error)
```
* 删除drawing。
* 我们会检查要删除的drawing是否为该uid所拥有，如果不属于删除会失败。 


### Drawing 类

```
GetID() string
```
* 取得 drawing ID。

```
Add(shape Shape) (err error)
```
* 添加新图形。

```
List() (shapes []Shape, err error)
```
* 列出所有图形。

```
Get(id ShapeID) (shape Shape, err error)
```
* 取出某个图形。

```
Set(id ShapeID, shape Shape) (err error)
```
* 修改某个图形。

```
SetZorder(id ShapeID, zorder string) (err error)
```
* 修改图形的图层。

```
Delete(id ShapeID) (err error)
```
* 删除某个图形。

```
Sync(shapes []ShapeID, changes []Shape) (err error)
```
* 同步drawing的修改。

### Shape 类型

```
type ShapeID = string

type Shape interface {
	GetID() ShapeID
}
```
* 考虑到 Drawing 类的 List 和 Get 返回的 Shape 实例，会被直接作为 RESTful API 的结果返回。所以 Shape json.Marshal 结果必须符合 API 层的预期。
* 考虑到 Drawing 类的 Add、Set、Sync 传入的 Shape 实例，会被直接写入 mongodb，所以 Shape bson.Marshal 结果必须符合 mongodb 的预期。

## 数据结构

基于 mongodb 来实现。

### drawing 表

| 字段名 | 含义 | 索引 | 类型 |
| ------- | ------ | ---------- | ------ |
| _id | DrawingID | 唯一索引 | string |
| uid | UserID | 索引 | UserID |
| shapes | Shape数组 | - | []ShapeID |

### shape 表

| 字段名 | 含义 | 索引 | 类型 |
| ------- | ------ | ---------- | ------ |
| dgid | DrawingID | - | string |
| spid | ShapeID | (dgid, spid) 联合唯一索引 | ShapeID |
| shape | Shape | - | json |

## 实现逻辑（算法）

创建新drawing (uid):
```
dgid = newObjectId()
db.drawing.insert({_id: dgid, uid: uid, shapes:[]})
return dgid
```

取得drawing的内容 (uid, dgid):
```
doc = db.drawing.findOne({_id: dgid, uid: uid})
shapes = []
foreach spid in doc.shapes {
    o = db.shape.find({dgid: dgid, spid: spid})
    shapes.push(o.shape)
}
return shapes
```

删除drawing (uid, dgid):
```
if db.drawing.remove({_id: dgid, uid: uid}) {
    db.shape.remove({dgid: dgid})
}
```

创建新shape (uid, dgid, shape):
```
if db.drawing.find({_id: dgid, uid: uid}) {
    db.shape.insert({dgid: dgid, spid: shape.id, shape: shape})
    db.drawing.update({$push: {shapes: shape.id}})
}
```

删除shape (uid, dgid, spid):
```
if db.drawing.find({_id: dgid, uid: uid}) {
    if db.drawing.update({$pull: {shapes: spid}}) {
        db.shape.remove({dgid: dgid, spid: spid})
    }
}
```
