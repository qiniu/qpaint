/*
Usage:
    <div type="BaseColorPicker" id="$id" onchange="onchange()" palette="$color1,$color2,...,$colorN"></div>
    color = black | black(黑色) | null | null(透明色) | null(transparent) | ...
*/
function BaseColorPicker(div) {
    let id = div.id
    let onchange = div.onchange
    let palette = div.getAttribute("palette")
    let colors = palette.split(",")
    let options = []
    for (let i in colors) {
        let color = colors[i]
        let n = color.length
        if (color.charAt(n-1) == ")") {
            let offset = color.indexOf("(")
            options.push(`<option value="` + color.substring(0, offset) + `">` + color.substring(offset+1, n-1) + `</option>`)
        } else {
            options.push(`<option value="` + color + `">` + color + `</option>`)
        }
    }
    div.outerHTML = `<select id="` + id + `">` + options.join("") + `</select>`
    let elem = document.getElementById(id)
    if (onchange) {
        elem.onchange = onchange
    }
}

qcontrols.register("BaseColorPicker", BaseColorPicker)
