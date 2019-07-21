# QPaint (by Qiniu.com)

## QPaint DOM

* TODO

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
