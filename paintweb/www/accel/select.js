class QShapeSelector {
    constructor() {
        this.started = false
        this.pt = {x: 0, y: 0}
        this.ptMove = {x: 0, y: 0}
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
        qview.drawing.style.cursor = "auto"
    }

    reset() {
        this.started = false
        invalidate(null)
    }

    onmousedown(event) {
        this.pt = this.ptMove = qview.getMousePos(event)
        this.started = true
        let ht = qview.doc.hitTest(this.pt)
        if (qview.selection != ht.hitShape) {
            qview.selection = ht.hitShape
            invalidate(null)
        }
    }
    onmousemove(event) {
        let pt = qview.getMousePos(event)
        if (this.started) {
            this.ptMove = pt
            invalidate(null)
        } else {
            let ht = qview.doc.hitTest(pt)
            if (ht.hitCode > 0) {
                qview.drawing.style.cursor = "move"
            } else {
                qview.drawing.style.cursor = "auto"
            }
        }
    }
    onmouseup(event) {
        if (this.started) {
            let selection = qview.selection
            if (selection != null) {
                let pt = qview.getMousePos(event)
                if (pt.x != this.pt.x || pt.y != this.pt.y) {
                    selection.move(qview.doc, pt.x - this.pt.x, pt.y - this.pt.y)
                }
            }
            this.reset()
        }
    }
    onkeydown(event) {
        switch (event.keyCode) {
        case 8:  // keyBackSpace
        case 46: // keyDelete
            qview.doc.deleteShape(qview.selection)
            qview.selection = null
        case 27: // keyEsc
            this.reset()
            break
        }
    }

    onpaint(ctx) {
        let selection = qview.selection
        if (selection != null) {
            let bound = selection.bound()
            if (this.started) {
                bound.x += this.ptMove.x - this.pt.x
                bound.y += this.ptMove.y - this.pt.y
            }
            ctx.lineWidth = 1
            ctx.strokeStyle = "gray"
            ctx.beginPath()
            ctx.setLineDash([5, 5])
            ctx.rect(bound.x, bound.y, bound.width, bound.height)
            ctx.stroke()
            ctx.setLineDash([])
        }
    }
}

qview.registerController("ShapeSelector", function() {
    return new QShapeSelector()
})
