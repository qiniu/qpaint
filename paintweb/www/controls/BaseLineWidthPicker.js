/*
Usage:
    <div type="BaseLineWidthPicker" id="$id" onchange="onchange()"></div>
*/
function BaseLineWidthPicker(div) {
    let id = div.id
    let onchange = div.onchange
    div.outerHTML = `<select id="` + id + `">
<option value="1">1</option>
<option value="3">3</option>
<option value="5">5</option>
<option value="7">7</option>
<option value="9">9</option>
<option value="11">11</option>
</select>`
    let elem = document.getElementById(id)
    if (onchange) {
        elem.onchange = onchange
    }
}

qcontrols.register("BaseLineWidthPicker", BaseLineWidthPicker)
