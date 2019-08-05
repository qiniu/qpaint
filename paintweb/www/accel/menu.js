// ----------------------------------------------------------

function unselectElementById(id, invalidKey) {
    if (id != invalidKey) {
        document.getElementById(id).removeAttribute("style")
    }
}

function selectElement(elem) {
    elem.style.borderColor = "blue"
}

function selectElementById(id, invalidKey) {
    if (id != invalidKey) {
        selectElement(document.getElementById(id))
    }
}

// ----------------------------------------------------------

function onClickCtrl(key) {
    unselectElementById(qview.currentKey, "ShapeSelector")
    let elem = document.getElementById(key)
    selectElement(elem)
    elem.blur()
    qview.invokeController(key)
}

function installControllers() {
    document.getElementById("menu").innerHTML = `
    <input type="button" id="PathCreator" value="Create Path" onclick="onClickCtrl('PathCreator')">
    <input type="button" id="FreePathCreator" value="Create FreePath" onclick="onClickCtrl('FreePathCreator')">
    <input type="button" id="LineCreator" value="Create Line" onclick="onClickCtrl('LineCreator')">
    <input type="button" id="RectCreator" value="Create Rect" onclick="onClickCtrl('RectCreator')">
    <input type="button" id="EllipseCreator" value="Create Ellipse" onclick="onClickCtrl('EllipseCreator')">
    <input type="button" id="CircleCreator" value="Create Circle" onclick="onClickCtrl('CircleCreator')">`

    onViewAdded(function(view) {
        view.invokeController("ShapeSelector")
        view.onControllerReset = function() {
            unselectElementById(view.currentKey, "ShapeSelector")
            view.invokeController("ShapeSelector")
        }
    })
    onCurrentViewChanged(function(old) {
        if (old) {
            unselectElementById(old.currentKey, "ShapeSelector")
        }
        selectElementById(qview.currentKey, "ShapeSelector")
    })
}

// ----------------------------------------------------------

function selection_setProp(key, val) {
    if (qview.selection != null) {
        qview.selection.setProp(qview.doc, key, val)
        invalidate(null)
    }
}

function onPropChanged(key) {
    let elem = document.getElementById(key)
    let val = elem.value
    elem.blur()
    qview.style[key] = val
    selection_setProp(key, val)
}

function onIntPropChanged(key) {
    let elem = document.getElementById(key)
    elem.blur()
    let val = parseInt(elem.value)
    if (val > 0) {
        qview.style[key] = val
        selection_setProp(key, val)
    }
}

function onSelectionChanged(old) {
    let selection = qview.selection
    if (selection != null) {
        let style = selection.style
        qview.style = style.clone()
        document.getElementById("lineWidth").value = style.lineWidth
        document.getElementById("lineColor").value = style.lineColor
        document.getElementById("fillColor").value = style.fillColor
    }
}

function installPropSelectors() {
    document.getElementById("menu").insertAdjacentHTML("afterend", `<br><div id="properties">
    <label for="lineWidth">LineWidth: </label>
    <div type="BaseLineWidthPicker" id="lineWidth" onchange="onIntPropChanged('lineWidth')"></div>&nbsp;
    <label for="lineColor">LineColor: </label>
    <div type="ColorPicker" id="lineColor" onchange="onPropChanged('lineColor')" palette="black,red,blue,green,yellow,gray"></div>&nbsp;
    <label for="fillColor">FillColor: </label>
    <div type="ColorPicker" id="fillColor" onchange="onPropChanged('fillColor')" palette="white,null(transparent),black,red,blue,green,yellow,gray"></div>
    </div>`)
    onViewAdded(function(view) {
        view.onSelectionChanged = function(old) {
            if (qview === view) {
                onSelectionChanged(old)
            }
        }
    })
}

// ----------------------------------------------------------

function installMousePos() {
    document.getElementById("properties").insertAdjacentHTML("beforeend", `&nbsp;<span id="mousepos"></span>`)

    let mousepos = document.getElementById("mousepos")
    onViewAdded(function(view) {
        let old = view.drawing.onmousemove
        view.drawing.onmousemove = function(event) {
            let pos = view.getMousePos(event)
            mousepos.innerText = "MousePos: " + pos.x + ", " + pos.y
            old(event)
        }
    })
}

// ----------------------------------------------------------

installControllers()
installPropSelectors()
installMousePos()

// ----------------------------------------------------------
