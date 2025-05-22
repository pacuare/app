import { language, query } from "../components/editor.js"
import { apiQuery } from "./api.js"
import { displayResults, downloadCurrentResults } from "./index/lib.js"
import "../components/accountsettings.js"

const runBtn = document.querySelector("#run-query")

async function runQuery() {
    const editor = document.querySelector("#query-editor")

    runBtn.setAttribute("disabled", "")
    runBtn.classList.add("disabled")
    
    await apiQuery(language(editor), query(editor)).then(displayResults);

    runBtn.removeAttribute("disabled")
    runBtn.classList.remove("disabled")
}

addEventListener('DOMContentLoaded', () => {
    runBtn.addEventListener("click", runQuery)

    document.querySelector("#open-docs").addEventListener("click", () => {
        document.querySelector("#docs-sidebar").classList.toggle("flex");
        ["#docs-sidebar", ".docs-icon-closed", ".docs-icon-open"].map(el =>
            document.querySelector(el).classList.toggle("hidden"));
    })

    document.querySelector("#exportCSV").addEventListener("click", downloadCurrentResults)

    if(location.search.includes("settings=") || location.search.includes("key=")) {
        document.querySelector("#openSettings").click()
        history.pushState(null, '', location.href.split('?')[0])
    }
})