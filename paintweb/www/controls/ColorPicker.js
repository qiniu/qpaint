function ColorPicker(div) {
    let id = div.id
    let onchange = div.onchange
    let palette = div.getAttribute("palette")
    let colors = palette.split(",")
    div.outerHTML = `<input type="button" id="` + id + `" value="` + colors[0] + `">`
    let elem = $("#" + id)
    elem.spectrum({
        showInitial: true,
        showInput: true,
        showButtons: true,
        preferredFormat: "hex6",
        change: function(color) {
            elem.val(color.toHexString())
            if (onchange) {
                onchange()
            }
        },
        move: function(color) {
            elem.val(color.toHexString())
        }
    })
}

qcontrols.register("ColorPicker", ColorPicker)
