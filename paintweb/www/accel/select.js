// ----------------------------------------------------------

class QShapeSelector {
    constructor(view) {
        this.started = false
        this.pt = {x: 0, y: 0}
        this.ptMove = {x: 0, y: 0}
        this.view = view
        let ctrl = this
        view.onmousedown = function(event) { ctrl.onmousedown(event) }
        view.onmousemove = function(event) { ctrl.onmousemove(event) }
        view.onmouseup = function(event) { ctrl.onmouseup(event) }
        view.onkeydown = function(event) { ctrl.onkeydown(event) }
    }
    stop() {
        let view = this.view
        view.onmousedown = null
        view.onmousemove = null
        view.onmouseup = null
        view.onkeydown = null
        view.drawing.style.cursor = "auto"
    }

    reset() {
        this.started = false
        this.view.invalidateRect(null)
    }

    onmousedown(event) {
        let view = this.view
        this.pt = this.ptMove = view.getMousePos(event)
        this.started = true
        let ht = view.doc.hitTest(this.pt)
        if (view.selection != ht.hitShape) {
            view.selection = ht.hitShape
            view.invalidateRect(null)
        }
    }
    onmousemove(event) {
        let view = this.view
        let pt = view.getMousePos(event)
        if (this.started) {
            this.ptMove = pt
            view.invalidateRect(null)
        } else {
            let ht = view.doc.hitTest(pt)
            if (ht.hitCode > 0) {
                view.drawing.style.cursor = "move"
            } else {
                view.drawing.style.cursor = "auto"
            }
        }
    }
    onmouseup(event) {
        if (this.started) {
            let view = this.view
            let selection = view.selection
            if (selection != null) {
                let pt = view.getMousePos(event)
                if (pt.x != this.pt.x || pt.y != this.pt.y) {
                    selection.move(view.doc, pt.x - this.pt.x, pt.y - this.pt.y)
                }
            }
            this.reset()
        }
    }
    onkeydown(event) {
        switch (event.keyCode) {
        case 8:  // keyBackSpace
        case 46: // keyDelete
            let view = this.view
            view.doc.deleteShape(view.selection)
            view.selection = null
        case 27: // keyEsc
            this.reset()
            break
        }
    }

    onpaint(ctx) {
        let selection = this.view.selection
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

onViewAdded(function(view) {
    view.registerController("ShapeSelector", function() {
        return new QShapeSelector(view)
    })
})

// ----------------------------------------------------------
