class QControls {
    constructor() {
        this.data = {}
    }

    register(kind, control) {
        this.data[kind] = control
    }

    init() {
        let divs = document.getElementsByTagName("div")
        let n = divs.length
        for (let i = n-1; i >= 0; i--) {
            let div = divs[i]
            let kind = div.getAttribute("type")
            if (kind != null) {
                let control = this.data[kind]
                if (control) {
                    control(div)
                }
            }
        }
    }
}

var qcontrols = new QControls()
