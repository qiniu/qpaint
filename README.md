# QPaint (by Qiniu.com)

[![Build Status](https://travis-ci.org/qiniu/qpaint.svg?branch=master)](https://travis-ci.org/qiniu/qpaint) [![GoDoc](https://godoc.org/github.com/qiniu/qpaint?status.svg)](https://godoc.org/github.com/qiniu/qpaint)

[![Qiniu Logo](http://open.qiniudn.com/logo.png)](http://www.qiniu.com/)

## 辅助界面元素

### 接口规格 (第 31 讲)

| 类型 | 属性 | 方法 | 事件 |
| ------------- | ---------- | ------ | ------------- |
| BaseLineWidthPicker | id: string<br>value: number | blur() | onchange(event: Event) |
| BaseColorPicker | id: string<br>value: Color<br>palette: string | blur() | onchange(event: Event) |
| ColorPicker | id: string<br>palette: string<br>value: Color | blur() | onchange(event: Event) |

## QPaint DOM (第 30 讲)

### 网络协议

在 29 讲的基础上，我们增加了以下能力：

| 功能 | 请求包 | 返回包 |
| ------------- | ---------- | ------------- |
| 同步变更 | POST /drawings/`<DrawingID>`/sync<br>Content-Type: `application/json`<br><br>{<br>&nbsp;&nbsp;&nbsp;&nbsp;"shapes": [`<ShapeID1>`, `<ShapeID2>`, `<ShapeID3>`, ...],<br>&nbsp;&nbsp;&nbsp;&nbsp;"changes": [<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;{<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"id": `<ShapeID>`,<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`<Shape>`<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;},<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;...<br>&nbsp;&nbsp;&nbsp;&nbsp;]<br>} | 200 OK |

## QPaint DOM (第 29 讲)

### 网络协议

| 功能 | 请求包 | 返回包 |
| ------------- | ---------- | ------------- |
| 创建新drawing | POST /drawings | 200 OK<br>Content-Type: `application/json`<br><br>{<br>&nbsp;&nbsp;&nbsp;&nbsp;"id": `<DrawingID>`<br>} |
| 获得drawing | GET /drawings/`<DrawingID>` | 200 OK<br>Content-Type: `application/json`<br><br>{<br>&nbsp;&nbsp;&nbsp;&nbsp;"shapes": [<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;{<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"id": `<ShapeID>`<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`<Shape>`<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;},<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;...<br>&nbsp;&nbsp;&nbsp;&nbsp;]<br>} |
| 删除drawing | DELETE /drawings/`<DrawingID>` | 200 OK |
| 创建新shape | POST /drawings/`<DrawingID>`/shapes<br>Content-Type: `application/json`<br><br>{<br>&nbsp;&nbsp;&nbsp;&nbsp;"id": `<ShapeID>`,<br>&nbsp;&nbsp;&nbsp;&nbsp;`<Shape>`<br>} | 200 OK |
| 取得shape | GET /drawings/`<DrawingID>`/shapes/`<ShapeID>` | 200 OK<br>Content-Type: `application/json`<br><br>{<br>&nbsp;&nbsp;&nbsp;&nbsp;`<Shape>`<br>} |
| 修改shape | POST /drawings/`<DrawingID>`/shapes/`<ShapeID>`<br>Content-Type: `application/json`<br><br>{<br>&nbsp;&nbsp;&nbsp;&nbsp;`<Shape>`<br>} | 200 OK |
| 修改shape的顺序 | POST /drawings/`<DrawingID>`/shapes/`<ShapeID>`<br>Content-Type: `application/json`<br><br>{<br>&nbsp;&nbsp;&nbsp;&nbsp;"zorder": `<ZorderOperation>`<br>} | 200 OK |
| 删除shape | DELETE /drawings/`<DrawingID>`/shapes/`<ShapeID>` | 200 OK |

其中 `<Shape>` 是这样的：

```
"path": {
    "points": [
        {"x": <X>, "y": <Y>},
        ...
    ],
    "close": <Boolean>,
    "style": <ShapeStyle>
}
```
或:
```
"line": {
    "pt1": {"x": <X>, "y": <Y>},
    "pt2": {"x": <X>, "y": <Y>},
    "style": <ShapeStyle>
}
```
或：
```
"rect": {
    "x": <X>,
    "y": <Y>,
    "width": <Width>,
    "height": <Height>,
    "style": <ShapeStyle>
}
```
或：
```
"ellipse": {
    "x": <X>,
    "y": <Y>,
    "radiusX": <RadiusX>,
    "radiusY": <RadiusY>,
    "style": <ShapeStyle>
}
```

其中 `<ShapeStyle>` 是这样的：

```
{
    "lineWidth": <Width>,  // 线宽
    "lineColor": <Color>,  // 线型颜色
    "fillColor": <Color>   // 填充色
}
```

其中 `<ZorderOperation>` 可能的值为：

* "top": 到最顶
* "bottom": 到最底
* "front": 往前一层
* "back": 往后一层

## QPaint Web (第 27 讲)

### Session-based Model

* [dom.js](https://github.com/qiniu/qpaint/blob/v27/paintweb/www/dom.js)

```TypeScript
interface Shape {
    onpaint(ctx: CanvasRenderingContext2D): void
    bound(): Rect
    hitTest(pt: Point): {hitCode: number, hitShape: Shape}
    setProp(key: string, val: any): void
    move(dx, dy: number): void
}
```

| 类型  | View | Controllers |
| ------------- | ---------- | ------------- |
| QPaintDoc | onpaint(ctx) | addShape(shape)<br>deleteShape(shape)<br>hitTest(pt) |
| QLine<br>QRect<br>QEllipse<br>QPath | onpaint(ctx) | new QLine(pt1, pt2, style)<br>new QRect(rect, style)<br>new QEllipse(x, y, radiusX, radiusY, style)<br>new QPath(points, close, style)<br>bound()<br>hitTest(pt)<br>setProp(key, val)<br>move(dx, dy) |
| QShapeStyle | new QShapeStyle(<br>&nbsp;&nbsp;lineWidth, lineColor, fillColor<br>) | setProp(key, val)<br>clone() |

### ViewModel

* [view.js](https://github.com/qiniu/qpaint/blob/v27/paintweb/www/view.js)
* [index.htm](https://github.com/qiniu/qpaint/blob/v27/paintweb/www/index.htm)

```TypeScript
interface Controller {
  stop(): void
  onpaint(ctx: CanvasRenderingContext2D): void
}
```

| 类型  | Model | View | Controllers |
| --------- | ------- | ---------- | ------------- |
| 数据 | doc: QPaintDoc | style: QShapeStyle<br>drawing: DOMElement | controllers: map[string]Controller |
| 方法 | - | invalidateRect(rect) | get currentKey()<br>get selection()<br>set selection(shape)<br>getMousePos(event)<br>registerController(name, ctrl)<br>invokeController(name)<br>stopController()<br>fireControllerReset() |
| 事件 | - | onpaint(ctx) | onmousedown(event)<br>onmousemove(event)<br>onmouseup(event)<br>ondblclick(event)<br>onkeydown(event)<br>onSelectionChanged(old)<br>onControllerReset() |

### Controllers

* Menu, PropSelectors, MousePosTracker: [accel/menu.js](https://github.com/qiniu/qpaint/blob/v27/paintweb/www/accel/menu.js)
* ShapeSelector: [accel/select.js](https://github.com/qiniu/qpaint/blob/v27/paintweb/www/accel/select.js)
* Create Path: [creator/path.js](https://github.com/qiniu/qpaint/blob/v27/paintweb/www/creator/path.js)
* Create FreePath: [creator/freepath.js](https://github.com/qiniu/qpaint/blob/v27/paintweb/www/creator/freepath.js)
* Create Line, Rect, Ellipse, Circle: [creator/rect.js](https://github.com/qiniu/qpaint/blob/v27/paintweb/www/creator/rect.js)

| 类型 | Event | Model | View |
| --- | --- | --- | --- |
| Menu | onControllerReset() | - | controllers: map[string]Controller<br>get currentKey()<br>invokeController(name) |
| PropSelectors | onSelectionChanged(old) | shape.style<br>shape.setProp(key, val)<br>style.clone() | style: QShapeStyle<br>get selection()<br>invalidateRect(rect) |
| MousePosTracker | onmousemove | - | getMousePos(event) |
| QShapeSelector | onmousedown(event)<br>onmousemove(event)<br>onmouseup(event)<br>onkeydown(event)<br>onpaint(ctx) | doc.deleteShape(shape)<br>doc.hitTest(pt)<br>shape.move(dx, dy)<br>shape.bound() | get selection()<br>set selection(shape)<br>getMousePos(event)<br>invalidateRect(rect)<br>registerController(name, ctrl) |
| QPathCreator | onmousedown(event)<br>onmousemove(event)<br>onmouseup(event)<br>ondblclick(event)<br>onkeydown(event)<br>onpaint(ctx) | new QPath(points, close, style)<br>doc.addShape(shape) | getMousePos(event)<br>invalidateRect(rect)<br>registerController(name, ctrl)<br>fireControllerReset() |
| QFreePathCreator | onmousedown(event)<br>onmousemove(event)<br>onmouseup(event)<br>onkeydown(event)<br>onpaint(ctx) | new QPath(points, close, style)<br>doc.addShape(shape) | getMousePos(event)<br>invalidateRect(rect)<br>registerController(name, ctrl)<br>fireControllerReset() |
| QRectCreator | onmousedown(event)<br>onmousemove(event)<br>onmouseup(event)<br>onkeydown(event)<br>onpaint(ctx) | new QLine(pt1, pt2, style)<br>new QRect(rect, style)<br>new QEllipse(x, y, radiusX, radiusY, style)<br>doc.addShape(shape) | getMousePos(event)<br>invalidateRect(rect)<br>registerController(name, ctrl)<br>fireControllerReset() |

### Change Notes

* https://github.com/qiniu/qpaint/compare/v26...v27

#### Session-based Model

| 类型  | View | Controllers | 修改说明 |
| ------- | ---------- | ------- | ------ |
| QPaintDoc | - | deleteShape(shape)<br>hitTest(pt) | - 增加 hitTest (确定鼠标点中哪个图形)、deleteShape (删除某个图形)，都用于 QShapeSelector |
| QLine<br>QRect<br>QEllipse<br>QPath | - | new QLine(pt1, pt2, style)<br>new QRect(rect, style)<br>new QEllipse(x, y, radiusX, radiusY, style)<br>new QPath(points, close, style)<br>bound()<br>hitTest(pt)<br>setProp(key, val)<br>move(dx, dy) | - 构造函数 style 参数由 QLineStyle 改为 QShapeStyle<br>- bound (求图形的外接矩形)、hitTest 用于选择图形<br>- setProp (修改图形样式的某个属性)<br>- move (移动图形) |
| QShapeStyle | new QShapeStyle(<br>&nbsp;&nbsp;lineWidth, lineColor, fillColor<br>) | setProp(key, val)<br>clone() | - QLineStyle 改名为 QShapeStyle<br>- 属性 width、color 改名为 lineWidth、lineColor<br>- 增加属性 fillColor (图形的填充色)<br>- 增加 setProp、clone (克隆图形样式) |

#### ViewModel

| 类型  | Model | View | Controllers | 修改说明 |
| --------- | ------- | ----- | ----- | ------ |
| 数据 | - | style: QShapeStyle | - | - 属性 properties 改名为 style
| 方法 | - | - | get selection()<br>set selection(shape)<br>fireControllerReset() | - 删除了 get lineStyle()，和 properties 统一为 style<br>- 增加了 selection 读写<br>- fireControllerReset，用于让创建图形的 Controller 完成或放弃图形创建时发出 onControllerReset 事件
| 事件 | - | - | onSelectionChanged(old)<br>onControllerReset() | - onSelectionChanged 在被选择的图形改变时发出<br>- onControllerReset (见 fireControllerReset 的说明) |

#### Controllers

| 类型 | Event | Model | View | 修改说明 |
| --- | --- | --- | --- | --- |
| Menu | onControllerReset() | - | - | - 引入了 v2 版本的切换 Controller 的范式，更接近现代的交互范式
| PropSelectors | onSelectionChanged(old) | shape.style<br>shape.setProp(key, val)<br>style.clone() | style: QShapeStyle<br>get selection()<br>invalidateRect(rect) | - 这个 Controller 要比上一版本的复杂很多：之前只是修改 view 的 properties (现在是 style) 属性，以便于创建图形时引用。现在是改变它时还会作用于 selection (被选中的图形)，改变它的样式；而且，在 selection 改变时，会自动更新界面以反映被选图形的样式 |
| MousePosTracker | - | - | - | - 无变化 |
| QShapeSelector | onmousedown(event)<br>onmousemove(event)<br>onmouseup(event)<br>onkeydown(event)<br>onpaint(ctx) | doc.deleteShape(shape)<br>doc.hitTest(pt)<br>shape.move(dx, dy)<br>shape.bound() | get selection()<br>set selection(shape)<br>getMousePos(event)<br>invalidateRect(rect)<br>registerController(name, ctrl) | - 完全新增的 Controller |
| QPathCreator | - | new QPath(points, close, style) | fireControllerReset() | - QPath 的 style 参数从 QLineStyle 变为 QShapeStyle<br>- 完成或放弃图形创建时发出 onControllerReset 事件 |
| QFreePathCreator | - | new QPath(points, close, style) | fireControllerReset() | - 同上 |
| QRectCreator | - | new QLine(pt1, pt2, style)<br>new QRect(rect, style)<br>new QEllipse(x, y, radiusX, radiusY, style) | fireControllerReset() | - 同上 |

## QPaint Web (第 26 讲)

### Session-based Model

* [dom.js](https://github.com/qiniu/qpaint/blob/v26/paintweb/www/dom.js)

```TypeScript
interface Shape {
    onpaint(ctx: CanvasRenderingContext2D): void
}
```

| 类型  | View | Controllers |
| ------------- | ---------- | ------------- |
| QPaintDoc | onpaint(ctx) | addShape(shape) |
| QLine<br>QRect<br>QEllipse<br>QPath | onpaint(ctx) | new QLine(pt1, pt2, style)<br>new QRect(rect, style)<br>new QEllipse(x, y, radiusX, radiusY, style)<br>new QPath(points, close, style) |
| QLineStyle | - | new QLineStyle(width, color) |

### ViewModel

* [view.js](https://github.com/qiniu/qpaint/blob/v26/paintweb/www/view.js)
* [index.htm](https://github.com/qiniu/qpaint/blob/v26/paintweb/www/index.htm)

```TypeScript
interface Controller {
  stop(): void
  onpaint(ctx: CanvasRenderingContext2D): void
}
```

| 类型  | Model | View | Controllers |
| --------- | ------- | ---------- | ------------- |
| 数据 | doc: QPaintDoc | properties: {<br>&nbsp;&nbsp;lineWidth: number<br>&nbsp;&nbsp;lineColor: string<br>}<br>drawing: DOMElement | controllers: map[string]Controller |
| 方法 | - | invalidateRect(rect) | get currentKey()<br>get lineStyle()<br>getMousePos(event)<br>registerController(name, ctrl)<br>invokeController(name)<br>stopController() |
| 事件 | - | onpaint(ctx) | onmousedown(event)<br>onmousemove(event)<br>onmouseup(event)<br>ondblclick(event)<br>onkeydown(event)<br> |

### Controllers

* Menu, PropSelectors, MousePosTracker: [accel/menu.js](https://github.com/qiniu/qpaint/blob/v26/paintweb/www/accel/menu.js)
* Create Path: [creator/path.js](https://github.com/qiniu/qpaint/blob/v26/paintweb/www/creator/path.js)
* Create FreePath: [creator/freepath.js](https://github.com/qiniu/qpaint/blob/v26/paintweb/www/creator/freepath.js)
* Create Line, Rect, Ellipse, Circle: [creator/rect.js](https://github.com/qiniu/qpaint/blob/v26/paintweb/www/creator/rect.js)

| 类型 | Event | Model | View |
| --- | --- | --- | --- |
| Menu | - | - | controllers: map[string]Controller<br>get currentKey()<br>invokeController(name) |
| PropSelectors | - | - | properties: {<br>&nbsp;&nbsp;lineWidth: number<br>&nbsp;&nbsp;lineColor: string<br>} |
| MousePosTracker | onmousemove | - | getMousePos(event) |
| QPathCreator | onmousedown(event)<br>onmousemove(event)<br>onmouseup(event)<br>ondblclick(event)<br>onkeydown(event)<br>onpaint(ctx) | new QPath(points, close, style)<br>doc.addShape(shape) | getMousePos(event)<br>invalidateRect(rect)<br>registerController(name, ctrl) |
| QFreePathCreator | onmousedown(event)<br>onmousemove(event)<br>onmouseup(event)<br>onkeydown(event)<br>onpaint(ctx) | new QPath(points, close, style)<br>doc.addShape(shape) | getMousePos(event)<br>invalidateRect(rect)<br>registerController(name, ctrl) |
| QRectCreator | onmousedown(event)<br>onmousemove(event)<br>onmouseup(event)<br>onkeydown(event)<br>onpaint(ctx) | new QLine(pt1, pt2, style)<br>new QRect(rect, style)<br>new QEllipse(x, y, radiusX, radiusY, style)<br>doc.addShape(shape) | getMousePos(event)<br>invalidateRect(rect)<br>registerController(name, ctrl) |
