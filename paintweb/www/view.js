class QPaintView {
    constructor() {
        this.properties = {
            lineWidth: 1,
            lineColor: "black"
        }
        this.controllers = {}
        this.currentKey = ""
        this.current = null
        this.onmousedown = null
        this.onmousemove = null
        this.onmouseup = null
        this.ondblclick = null
        this.onkeypress = null
        this.onCurrentChanged = null
        let drawing = document.getElementById("drawing")
        let view = this
        drawing.onmousedown = function(event) {
            event.preventDefault()
            if (view.onmousedown != null) {
                view.onmousedown(event)
            }
        }
        drawing.onmousemove = function(event) {
            if (view.onmousemove != null) {
                view.onmousemove(event)
            }
        }
        drawing.onmouseup = function(event) {
            if (view.onmouseup != null) {
                view.onmouseup(event)
            }
        }
        drawing.ondblclick = function(event) {
            event.preventDefault()
            if (view.ondblclick != null) {
                view.ondblclick(event)
            }
        }
        document.onkeydown = function(event) {
            event.preventDefault()
            if (view.onkeydown != null) {
                view.onkeydown(event)
            }
        }
        this.drawing = drawing
        this.bound = drawing.getBoundingClientRect()
    }

    onpaint(ctx) {
        if (this.current != null) {
            this.current.onpaint(ctx)
        }
    }

    getMousePos(event) {
        return {
            x: event.offsetX,
            y: event.offsetY
        }
    }

    invalidateRect(reserved) {
        let ctx = this.drawing.getContext("2d")
        ctx.clearRect(0, 0, this.bound.width, this.bound.height)
        this.onpaint(ctx)
    }

    registerController(name, controller) {
        if (name in this.controllers) {
            alert("Controller exists: " + name)
        } else {
            this.controllers[name] = controller
        }
    }
    invokeController(name) {
        this.stopController()
        if (name in this.controllers) {
            let controller = this.controllers[name]
            this._setCurrent(name, controller())
        }
    }
    stopController() {
        if (this.current != null) {
            this.current.stop()
            this._setCurrent("", null)
        }
    }

    _setCurrent(name, ctrl) {
        this.current = ctrl
        this.currentKey = name
        if (this.onCurrentChanged) {
            this.onCurrentChanged()
        }
    }
}

var qview = new QPaintView()

function invalidate(reserved) {
    qview.invalidateRect(null)
}
