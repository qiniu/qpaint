class QFreePathCreator {
    constructor() {
        this.points = []
        this.fromPos = {x: 0, y: 0}
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
        this.points = []
        this.started = false
        invalidate(null)
    }

    onmousedown(event) {
        this.fromPos = qview.getMousePos(event)
        this.started = true
    }
    onmousemove(event) {
        if (this.started) {
            this.points.push(qview.getMousePos(event))
            invalidate(null)
        }
    }
    onmouseup(event) {
        if (this.started) {
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
            ctx.lineWidth = qview.properties.lineWidth
            ctx.beginPath()
            ctx.moveTo(this.fromPos.x, this.fromPos.y)
            for (let i in this.points) {
                ctx.lineTo(this.points[i].x, this.points[i].y)
            }
            ctx.stroke()
        }
    }
}

qview.registerController("FreePathCreator", function() {
    return new QFreePathCreator()
})
