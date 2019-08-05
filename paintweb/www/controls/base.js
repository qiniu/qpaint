class QControls {
    constructor() {
        this.data = {}
    }

    register(type, control) {
        this.data[type] = control
    }

    init() {
        let divs = document.getElementsByTagName("div")
        let n = divs.length
        for (let i = n-1; i >= 0; i--) {
            let div = divs[i]
            let type = div.getAttribute("type")
            if (type != null) {
                let control = this.data[type]
                if (control) {
                    control(div)
                }
            }
        }
    }
}

var qcontrols = new QControls()
