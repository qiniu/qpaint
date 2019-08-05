class QPathCreator {
    constructor(view, close) {
        this.points = []
        this.close = close
        this.fromPos = this.toPos = {x: 0, y: 0}
        this.started = false
        this.view = view
        let ctrl = this
        view.onmousedown = function(event) { ctrl.onmousedown(event) }
        view.onmousemove = function(event) { ctrl.onmousemove(event) }
        view.ondblclick = function(event) { ctrl.ondblclick(event) }
        view.onkeydown = function(event) { ctrl.onkeydown(event) }
    }
    stop() {
        let view = this.view
        view.onmousedown = null
        view.onmousemove = null
        view.ondblclick = null
        view.onkeydown = null
    }

    reset() {
        this.points = []
        this.started = false
        let view = this.view
        view.invalidateRect(null)
        view.fireControllerReset()
    }
    buildShape() {
        let points = [{x: this.fromPos.x, y: this.fromPos.y}]
        for (let i in this.points) {
            points.push(this.points[i])
        }
        return new QPath(points, this.close, defaultStyle.clone())
    }

    onmousedown(event) {
        let view = this.view
        this.toPos = view.getMousePos(event)
        if (this.started) {
            this.points.push(this.toPos)
        } else {
            this.fromPos = this.toPos
            this.started = true
        }
        view.invalidateRect(null)
    }
    onmousemove(event) {
        if (this.started) {
            let view = this.view
            this.toPos = view.getMousePos(event)
            view.invalidateRect(null)
        }
    }
    ondblclick(event) {
        if (this.started) {
            this.view.doc.addShape(this.buildShape())
            this.reset()
        }
    }
    onkeydown(event) {
        switch (event.keyCode) {
        case 13: // keyEnter
            let n = this.points.length
            if (n == 0 || this.points[n-1] !== this.toPos) {
                this.points.push(this.toPos)
            }
            this.ondblclick(event)
            break
        case 27: // keyEsc
            this.reset()
        }
    }

    onpaint(ctx) {
        if (this.started) {
            let props = defaultStyle
            ctx.lineWidth = props.lineWidth
            ctx.strokeStyle = props.lineColor
            ctx.beginPath()
            ctx.moveTo(this.fromPos.x, this.fromPos.y)
            for (let i in this.points) {
                ctx.lineTo(this.points[i].x, this.points[i].y)
            }
            ctx.lineTo(this.toPos.x, this.toPos.y)
            if (this.close) {
                ctx.closePath()
            }
            ctx.stroke()
        }
    }
}

onViewAdded(function(view) {
    view.registerController("PathCreator", function() {
        return new QPathCreator(view, false)
    })
})
