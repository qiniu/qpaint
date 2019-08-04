function ColorPicker(div) {
    let id = div.id
    let onchange = div.onchange
    let palette = div.getAttribute("palette")
    let colors = palette.split(",")
    let style = `" style="border-color:` + colors[0] + `;background:` + colors[0]
    div.outerHTML = `
    <input type="button" id="` + id + `" value="` + colors[0] + style + `"></input>`
    let elem = document.getElementById(id)
    if (onchange) {
        elem.onchange = onchange
    }
    elem.onclick = function() {
        alert("click ColorPicker")
        elem.blur()
    }
}

qcontrols.register("ColorPicker", ColorPicker)
