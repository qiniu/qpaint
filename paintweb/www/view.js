class QPaintView {
    constructor() {
        this.properties = {
            lineWidth: 1,
            lineColor: "black"
        }
        this.controllers = {}
        this._currentKey = ""
        this._current = null
        this.onmousedown = null
        this.onmousemove = null
        this.onmouseup = null
        this.ondblclick = null
        this.onkeypress = null
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
            switch (event.keyCode) {
            case 9: case 13: case 27:
                event.preventDefault()
            }
            if (view.onkeydown != null) {
                view.onkeydown(event)
            }
        }
        this.drawing = drawing
        this.doc = new QPaintDoc()
    }

    get currentKey() {
        return this._currentKey
    }
    get lineStyle() {
        let props = this.properties
        return new QLineStyle(props.lineWidth, props.lineColor)
    }

    onpaint(ctx) {
        this.doc.onpaint(ctx)
        if (this._current != null) {
            this._current.onpaint(ctx)
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
        let bound = this.drawing.getBoundingClientRect()
        ctx.clearRect(0, 0, bound.width, bound.height)
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
        if (this._current != null) {
            this._current.stop()
            this._setCurrent("", null)
        }
    }

    _setCurrent(name, ctrl) {
        this._current = ctrl
        this._currentKey = name
    }
}

var qview = new QPaintView()

function invalidate(reserved) {
    qview.invalidateRect(null)
}
