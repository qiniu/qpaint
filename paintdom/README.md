QPaint DOM API
========

为了简化，我们不引入多用户，只是引入多 drawing。

todo: 鉴权与安全

## 创建新drawing

请求包：

```
POST /drawings
```

返回包：

```
200 OK
Content-Type: application/json

{
    "id": <DrawingID>
}
```

## 获得drawing

请求包：

```
GET /drawings/<DrawingID>
```

返回包：

```
200 OK
Content-Type: application/json

{
    "shapes": [
        {
            "id": <ShapeID>
            <Shape>
        },
        ...
    ]
}
```

## 删除drawing

请求包：

```
DELETE /drawings/<DrawingID>
```

返回包：

```
200 OK
```

## 创建新shape

请求包：

```
POST /drawings/<DrawingID>/shapes
Content-Type: application/json

{
    "id": <ShapeID>,
    <Shape>
}
```

返回包：

```
200 OK
```

这里 `Shape` 是这样的：

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

或者：

```
"line": {
    "pt1": {"x": <X>, "y": <Y>},
    "pt2": {"x": <X>, "y": <Y>},
    "style": <ShapeStyle>
}
```

或者：

```
"rect": {
    "x": <X>,
    "y": <Y>,
    "width": <Width>,
    "height": <Height>,
    "style": <ShapeStyle>
}
```

或者：

```
"ellipse": {
    "x": <X>,
    "y": <Y>,
    "radiusX": <RadiusX>,
    "radiusY": <RadiusY>,
    "style": <ShapeStyle>
}
```

这里 `ShapeStyle` 是这样的：

```
{
    "lineWidth": <Width>,  // 线宽
    "lineColor": <Color>,  // 线型颜色
    "fillColor": <Color>,  // 填充色
}
```

## 取得shape

请求包：

```
GET /drawings/<DrawingID>/shapes/<ShapeID>
```

返回包：

```
200 OK
Content-Type: application/json

{
    <Shape>
}
```

## 修改shape

请求包：

```
POST /drawings/<DrawingID>/shapes/<ShapeID>
Content-Type: application/json

{
    <Shape>
}
```

返回包：

```
200 OK
```

## 修改shape的顺序

请求包：

```
POST /drawings/<DrawingID>/shapes/<ShapeID>
Content-Type: application/json

{
    "zorder": <ZorderOperation>
}
```

返回包：

```
200 OK
```

这里 `ZorderOperation` 可能是：

* "top": 到最顶
* "bottom": 到最底
* "front": 往前一层
* "back": 往后一层

## 删除shape

请求包：

```
DELETE /drawings/<DrawingID>/shapes/<ShapeID>
```

返回包：

```
200 OK
```
