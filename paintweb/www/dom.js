
// d = |(y2-y1)x0 + (x1-x2)y0 + (x2y1 - x1y2)| / sqrt[(x1-x2)^2+(y1-y2)^2]
//
function hitLine(pt, pt1, pt2, width) {
    if ((pt1.x - pt.x) * (pt.x - pt2.x) < 0) {
        return false
    }
    if ((pt1.y - pt.y) * (pt.y - pt2.y) < 0) {
        return false
    }
    let dy = pt2.y - pt1.y
    let dx = pt2.x - pt1.x
    let d12 = Math.sqrt(dx*dx + dy*dy)
    if (d12 < 0.1) {
        return false
    }
    let d = Math.abs(dy*pt.x - dx*pt.y + pt2.x*pt1.y - pt1.x*pt2.y) / d12 - 2
    return width >= d*2
}

function hitRect(pt, r) {
    if ((r.x + r.width - pt.x) * (pt.x - r.x) < 0) {
        return false
    }
    if ((r.y + r.height - pt.y) * (pt.y - r.y) < 0) {
        return false
    }
    return true
}

function normalizeRect(rect) {
    let x = rect.pt1.x
    let y = rect.pt1.y
    let width = rect.pt2.x - x
    let height = rect.pt2.y - y
    if (width < 0) {
        x = rect.pt2.x
        width = -width
    }
    if (height < 0) {
        y = rect.pt2.y
        height = -height
    }
    return {x: x, y: y, width: width, height: height}
}

function fill(ctx, fillColor) {
    if (fillColor != "null") {
        ctx.fillStyle = fillColor
        ctx.fill()
    }
}

function deleteItem(array, item) {
    let index = array.indexOf(item)
    if (index !== -1) {
        array.splice(index, 1)
    }
}

class QShapeStyle {
    constructor(lineWidth, lineColor, fillColor) {
        this.lineWidth = lineWidth
        this.lineColor = lineColor
        this.fillColor = fillColor
    }
    setProp(key, val) {
        this[key] = val
    }
    clone() {
        return new QShapeStyle(this.lineWidth, this.lineColor, this.fillColor)
    }
}

class QLine {
    constructor(point1, point2, style) {
        this.pt1 = point1
        this.pt2 = point2
        this.style = style
    }

    bound() {
        return normalizeRect(this)
    }
    hitTest(pt) {
        if (hitLine(pt, this.pt1, this.pt2, this.style.lineWidth)) {
            return {hitCode: 1, hitShape: this}
        }
        return {hitCode: 0, hitShape: null}
    }
    move(dx, dy) {
        this.pt1.x += dx
        this.pt1.y += dy
        this.pt2.x += dx
        this.pt2.y += dy
    }
    setProp(key, val) {
        this.style.setProp(key, val)
    }

    onpaint(ctx) {
        let style = this.style
        ctx.lineWidth = style.lineWidth
        ctx.strokeStyle = style.lineColor
        ctx.beginPath()
        ctx.moveTo(this.pt1.x, this.pt1.y)
        ctx.lineTo(this.pt2.x, this.pt2.y)
        ctx.stroke()
    }
}

class QRect {
    constructor(r, style) {
        this.x = r.x
        this.y = r.y
        this.width = r.width
        this.height = r.height
        this.style = style
    }

    bound() {
        return {x: this.x, y: this.y, width: this.width, height: this.height}
    }
    hitTest(pt) {
        if (hitRect(pt, this)) {
            return {hitCode: 1, hitShape: this}
        }
        return {hitCode: 0, hitShape: null}
    }
    move(dx, dy) {
        this.x += dx
        this.y += dy
    }
    setProp(key, val) {
        this.style.setProp(key, val)
    }

    onpaint(ctx) {
        let style = this.style
        ctx.lineWidth = style.lineWidth
        ctx.strokeStyle = style.lineColor
        ctx.beginPath()
        ctx.rect(this.x, this.y, this.width, this.height)
        fill(ctx, style.fillColor)
        ctx.stroke()
    }
}

class QEllipse {
    constructor(x, y, radiusX, radiusY, style) {
        this.x = x
        this.y = y
        this.radiusX = radiusX
        this.radiusY = radiusY
        this.style = style
    }

    bound() {
        return {
            x: this.x - this.radiusX,
            y: this.y - this.radiusY,
            width: this.radiusX * 2,
            height: this.radiusY * 2
        }
    }
    hitTest(pt) {
        let dx = pt.x - this.x
        let dy = pt.y - this.y
        let a = this.radiusX
        let b = this.radiusY
        if (dx*dx/a/a + dy*dy/b/b <= 1) {
            return {hitCode: 1, hitShape: this}
        }
        return {hitCode: 0, hitShape: null}
    }
    move(dx, dy) {
        this.x += dx
        this.y += dy
    }
    setProp(key, val) {
        this.style.setProp(key, val)
    }

    onpaint(ctx) {
        let style = this.style
        ctx.lineWidth = style.lineWidth
        ctx.strokeStyle = style.lineColor
        ctx.beginPath()
        ctx.ellipse(this.x, this.y, this.radiusX, this.radiusY, 0, 0, 2 * Math.PI)
        fill(ctx, style.fillColor)
        ctx.stroke()
    }
}

class QPath {
    constructor(points, close, style) {
        this.points = points
        this.close = close
        this.style = style
    }

    bound() {
        let points = this.points
        let n = points.length
        if (n < 1) {
            return
        }
        let x1 = points[0].x
        let y1 = points[0].y
        let x2 = x1
        let y2 = y1
        for (let i = 1; i < n; i++) {
            let tx = points[i].x
            let ty = points[i].y
            if (tx < x1) {
                x1 = tx
            } else if (tx > x2) {
                x2 = tx
            }
            if (ty < y1) {
                y1 = ty
            } else if (ty > y2) {
                y2 = ty
            }
        }        
        return {x: x1, y: y1, width: x2 - x1, height: y2 - y1}
    }
    hitTest(pt) {
        if (hitRect(pt, this.bound())) {
            let points = this.points
            let n = points.length
            if (n > 1) {
                let lineWidth = this.style.lineWidth
                for (let i = 1; i < n; i++) {
                    if (hitLine(pt, points[i-1], points[i], lineWidth)) {
                        return {hitCode: 1, hitShape: this}
                    }
                }
            }
        }
        return {hitCode: 0, hitShape: null}
    }
    move(dx, dy) {
        let points = this.points
        for (let i in points) {
            points[i].x += dx
            points[i].y += dy
        }
    }
    setProp(key, val) {
        this.style.setProp(key, val)
    }

    onpaint(ctx) {
        let points = this.points
        let n = points.length
        if (n < 1) {
            return
        }
        let style = this.style
        ctx.lineWidth = style.lineWidth
        ctx.strokeStyle = style.lineColor
        ctx.beginPath()
        ctx.moveTo(points[0].x, points[0].y)
        for (let i = 1; i < n; i++) {
            ctx.lineTo(points[i].x, points[i].y)
        }
        if (this.close) {
            ctx.closePath()
        }
        ctx.stroke()
    }
}

class QPaintDoc {
    constructor() {
        this.shapes = []
    }

    addShape(shape) {
        if (shape != null) {
            this.shapes.push(shape)
        }
    }
    deleteShape(shape) {
        deleteItem(this.shapes, shape)
    }
    hitTest(pt) {
        let shapes = this.shapes
        let n = shapes.length
        for (let i = n-1; i >= 0; i--) {
            let ret = shapes[i].hitTest(pt)
            if (ret.hitCode > 0) {
                return ret
            }
        }
        return {hitCode: 0, hitShape: null}
    }

    onpaint(ctx) {
        let shapes = this.shapes
        for (let i in shapes) {
            shapes[i].onpaint(ctx)
        }
    }
}
