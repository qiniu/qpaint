
function normalizeRect(rect) {
    let x = rect.p1.x
    let y = rect.p1.y
    let width = rect.p2.x - x
    let height = rect.p2.y - y
    if (width < 0) {
        x = rect.p2.x
        width = -width
    }
    if (height < 0) {
        y = rect.p2.y
        height = -height
    }
    return {x: x, y: y, width: width, height: height}
}

class QRectCreator {
    constructor(shapeType) {
        this.shapeType = shapeType
        this.rect = {
            p1: {x: 0, y: 0},
            p2: {x: 0, y: 0}
        }
        this.started = false
        let ctrl = this
        qview.onmousedown = function(event) { ctrl.onmousedown(event) }
        qview.onmousemove = function(event) { ctrl.onmousemove(event) }
        qview.onmouseup = function(event) { ctrl.onmouseup(event) }
        qview.onkeydown = function(event) { ctrl.onkeydown(event) }
    }
    stop() {
        qview.onmousedown = null
        qview.onmousemove = null
        qview.onmouseup = null
        qview.onkeydown = null
    }

    reset() {
        this.started = false
        invalidate(this.rect)
    }
    buildShape() {
        let rect = this.rect
        let r = normalizeRect(rect)
        switch (this.shapeType) {
        case "line":
            return new QLine(rect.p1, rect.p2, qview.lineStyle)
        case "rect":
            return new QRect(r, qview.lineStyle)
        case "ellipse":
            let rx = r.width / 2
            let ry = r.height / 2
            return new QEllipse(r.x + rx, r.y + ry, rx, ry, qview.lineStyle)
        case "circle":
            let rc = Math.sqrt(r.width * r.width + r.height * r.height)
            return new QEllipse(rect.p1.x, rect.p1.y, rc, rc, qview.lineStyle)
        default:
            alert("unknown shapeType: " + this.shapeType)
            return null
        }
    }

    onmousedown(event) {
        this.rect.p1 = qview.getMousePos(event)
        this.started = true
    }
    onmousemove(event) {
        if (this.started) {
            this.rect.p2 = qview.getMousePos(event)
            invalidate(this.rect)
        }
    }
    onmouseup(event) {
        if (this.started) {
            this.rect.p2 = qview.getMousePos(event)
            qview.doc.addShape(this.buildShape())
            this.reset()
        }
    }
    onkeydown(event) {
        if (event.keyCode == 27) { // keyEsc
            this.reset()
        }
    }

    onpaint(ctx) {
        if (this.started) {
            this.buildShape().onpaint(ctx)
        }
    }
}

qview.registerController("LineCreator", function() {
    return new QRectCreator("line")
})

qview.registerController("RectCreator", function() {
    return new QRectCreator("rect")
})

qview.registerController("EllipseCreator", function() {
    return new QRectCreator("ellipse")
})

qview.registerController("CircleCreator", function() {
    return new QRectCreator("circle")
})
