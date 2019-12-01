/*
Usage:
    <div type="ColorPicker" id="$id" onchange="onchange()"></div>
*/
function ColorPicker(div) {
    let id = div.id
    let onchange = div.onchange
    let palette = div.getAttribute("palette")
    let colors = palette.split(",")
    let value = colors[0]
    div.outerHTML = `<input type="button" id="` + id + `" value="` + value + `">`
    let elem = $("#" + id)
    elem.spectrum({
        showInitial: true,
        showInput: true,
        showButtons: true,
        preferredFormat: "hex6"
    })
    if (onchange) {
        elem.change(onchange)
    }
    Object.defineProperty(document.getElementById(id), "value", {
        get() {
            return value
        },
        set(x) {
            if (this.busy) {
                return
            }
            value = x
            this.busy = true
            elem.spectrum("set", value)
            this.busy = false
        }
    })
}

qcontrols.register("ColorPicker", ColorPicker)
