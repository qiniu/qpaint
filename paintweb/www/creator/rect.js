class QRectCreator {
    constructor(shapeType) {
        this.shapeType = shapeType
        this.rect = {
            pt1: {x: 0, y: 0},
            pt2: {x: 0, y: 0}
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
        qview.fireControllerReset()
    }
    buildShape() {
        let rect = this.rect
        let r = normalizeRect(rect)
        let style = qview.style.clone()
        switch (this.shapeType) {
        case "line":
            return new QLine(rect.pt1, rect.pt2, style)
        case "rect":
            return new QRect(r, style)
        case "ellipse":
            let rx = r.width / 2
            let ry = r.height / 2
            return new QEllipse(r.x + rx, r.y + ry, rx, ry, style)
        case "circle":
            let rc = Math.sqrt(r.width * r.width + r.height * r.height)
            return new QEllipse(rect.pt1.x, rect.pt1.y, rc, rc, style)
        default:
            alert("unknown shapeType: " + this.shapeType)
            return null
        }
    }

    onmousedown(event) {
        this.rect.pt1 = qview.getMousePos(event)
        this.started = true
    }
    onmousemove(event) {
        if (this.started) {
            this.rect.pt2 = qview.getMousePos(event)
            invalidate(this.rect)
        }
    }
    onmouseup(event) {
        if (this.started) {
            this.rect.pt2 = qview.getMousePos(event)
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
