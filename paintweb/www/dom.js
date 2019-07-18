
// d = |(y2-y1)x0 + (x1-x2)y0 + (x2y1 - x1y2)| / sqrt[(x1-x2)^2+(y1-y2)^2]
//
function hitOn(pt, pt1, pt2, width) {
    let dy = pt2.y - pt1.y
    let dx = pt2.x - pt1.x
    let d12 = Math.sqrt(dx*dx + dy*dy)
    if (d12 < 0.1) {
        let dy0 = pt1.y - pt.y
        let dx0 = pt1.x - pt.x
        let d01 = dx0*dx0 + dy0*dy0
        let w1 = width/2 + 1
        return w1*w1 >= d01
    }
    let d = Math.abs(dy*pt.x - dx*pt.y + pt2.x*pt1.y - pt1.x*pt2.y) / d12 - 1
    return width >= d*2 
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
        if (hitOn(pt, this.pt1, this.pt2, this.style.lineWidth)) {
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
        let y2 = y2
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
        return {hitCode: 0, hitShape: null}
    }
    move(dx, dy) {
        let points = this.points
        for (i in points) {
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
