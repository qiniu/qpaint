
class QBaseColorPicker {
    constructor(div) {
        let id = div.id
        let onchange = div.onchange
        div.outerHTML = `<select id="` + id + `">
<option value="black">black</option>
<option value="red">red</option>
<option value="blue">blue</option>
<option value="green">green</option>
<option value="yellow">yellow</option>
<option value="gray">gray</option>
</select>`
        let elem = document.getElementById(id)
        if (onchange) {
            elem.onchange = onchange
        }
    }
}

qcontrols.register("BaseColorPicker", function(div) {
    return new QBaseColorPicker(div)
})
