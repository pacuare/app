import "../components/accountsettings.js"
import Alpine from "alpinejs"
import './results.js'

Alpine.data('index', () => ({
    docsOpen: false,

    openDocs: {
        ['@click']() {
            this.docsOpen = !this.docsOpen
        }
    },

    exportCSV: {
        ['@click']() {
            this.$store.results.download();
        }
    },

    init() {
        if(location.search.includes("settings=") || location.search.includes("key=")) {
            document.querySelector("#openSettings").click()
            history.pushState(null, '', location.href.split('?')[0])
        }
    }
}))

Alpine.start()