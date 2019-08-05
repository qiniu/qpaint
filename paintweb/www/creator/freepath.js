class QFreePathCreator {
    constructor(view) {
        this.points = []
        this.fromPos = {x: 0, y: 0}
        this.started = false
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
        this.fromPos = this.view.getMousePos(event)
        this.started = true
    }
    onmousemove(event) {
        if (this.started) {
            let view = this.view
            this.points.push(view.getMousePos(event))
            view.invalidateRect(null)
        }
    }
    onmouseup(event) {
        if (this.started) {
            this.view.doc.addShape(this.buildShape())
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
            let props = defaultStyle
            ctx.lineWidth = props.lineWidth
            ctx.strokeStyle = props.lineColor
            ctx.beginPath()
            ctx.moveTo(this.fromPos.x, this.fromPos.y)
            for (let i in this.points) {
                ctx.lineTo(this.points[i].x, this.points[i].y)
            }
            ctx.stroke()
        }
    }
}

onViewAdded(function(view) {
    view.registerController("FreePathCreator", function() {
        return new QFreePathCreator(view)
    })
})
