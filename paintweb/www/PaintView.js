function newPaintView(drawingID) {
    let view = new QPaintView(drawingID)
    registerFreePathCreator(view)
    registerPathCreator(view)
    registerLineCreator(view)
    registerRectCreator(view)
    registerEllipseCreator(view)
    registerCircleCreator(view)
    registerShapeSelector(view)
    return view
}

function initPaintView(drawingID) {
    let view = newPaintView(drawingID)
    setCurrentView(view)
}

function PaintView(div) {
    let id = div.id
    let width = div.getAttribute("width")
    let height = div.getAttribute("height")
    div.outerHTML = `<canvas id="` + id + `" width="` + width + `" height="` + height + `">你的浏览器不支持Canvas！</canvas>`
    initPaintView(id)
}

qcontrols.register("PaintView", PaintView)
