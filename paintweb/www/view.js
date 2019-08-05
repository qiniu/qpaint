// ----------------------------------------------------------

var qview = null
var _onCurrentViewChangeds = []

function onCurrentViewChanged(handle) {
    _onCurrentViewChangeds.push(handle)
}

function setCurrentView(view) {
    let old = qview
    qview = view
    for (let i in _onCurrentViewChangeds) {
        let handle = _onCurrentViewChangeds[i]
        handle(old)
    }
}

function invalidate(reserved) {
    qview.invalidateRect(reserved)
}

// ----------------------------------------------------------

var _onViewAddeds = []

function onViewAdded(handle) {
    _onViewAddeds.push(handle)
}

function fireViewAdded(view) {
    for (let i in _onViewAddeds) {
        let handle = _onViewAddeds[i]
        handle(view)
    }
}

// ----------------------------------------------------------

class QPaintView {
    constructor(drawingID) {
        this.controllers = {}
        this._currentKey = ""
        this._current = null
        this._selection = null
        this.onmousedown = null
        this.onmousemove = null
        this.onmouseup = null
        this.ondblclick = null
        this.onkeydown = null
        this.onSelectionChanged = null
        this.onControllerReset = null
        let drawing = document.getElementById(drawingID)
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
        drawing.onmouseenter = function(event) {
            setCurrentView(view)
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
        window_onhashchange(function(event) {
            view.doc.reload()
            view.invalidateRect(null)
        })
        this.drawing = drawing
        this.doc = new QPaintDoc()
        this.doc.onload = function() {
            view.invalidateRect(null)
        }
        this.doc.init()
    }

    get currentKey() {
        return this._currentKey
    }
    get selection() {
        return this._selection
    }
    set selection(shape) {
        let old = this._selection
        if (old != shape) {
            this._selection = shape
            if (this.onSelectionChanged != null) {
                this.onSelectionChanged(old)
            }
        }
    }

    getMousePos(event) {
        return {
            x: event.offsetX,
            y: event.offsetY
        }
    }
 
    onpaint(ctx) {
        this.doc.onpaint(ctx)
        if (this._current != null) {
            this._current.onpaint(ctx)
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
    fireControllerReset() {
        if (this.onControllerReset != null) {
            this.onControllerReset()
        }
    }

    _setCurrent(name, ctrl) {
        this._current = ctrl
        this._currentKey = name
    }
}

// ----------------------------------------------------------
