import { downloadResults } from "./download.js";

let currentResults = []

export function displayResults(res) {
    currentResults = res;
    const table = document.querySelector("#resultsTable")
    const dlBtn = document.querySelector("#exportCSV")
    const err = document.querySelector("#resultError")
    if('error' in res) {
        table.classList.add("hidden")
        dlBtn.classList.add("hidden")
        dlBtn.classList.remove("inline-flex")
        err.classList.remove("hidden")
        err.textContent = res.error
    } else {
        if(res.length == 0) {
            table.classList.add("hidden")
            dlBtn.classList.add("hidden")
            dlBtn.classList.remove("inline-flex")
            err.classList.remove("hidden")
            err.textContent = "No rows in response"
            return
        }
        table.innerHTML = ""
        table.classList.remove("hidden")
        dlBtn.classList.remove("hidden")
        dlBtn.classList.add("inline-flex")
        err.classList.add("hidden")
        const fields = Object.keys(res[0]);
        const th = document.createElement("thead");
        const thr = document.createElement("tr");
        
        for(const field of fields) {
            const thd = document.createElement("th");
            thd.textContent = field;
            thr.appendChild(thd);
        }
        th.appendChild(thr);
        table.appendChild(th);

        const tb = document.createElement("tbody");
        tb.className = "[&_tr:last-child]:border-0"
        for(const row of res) {
            const tr = document.createElement("tr");
            for(const value of Object.values(row)) {
                const td = document.createElement("td");
                td.textContent = value?.toString() ?? "null";
                tr.appendChild(td);
            }
            tb.appendChild(tr);
        }
        table.appendChild(tb);
    }
}

export function downloadCurrentResults() {
    downloadResults(currentResults)
}