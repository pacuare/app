import { language, query } from "../components/editor.js"

const runBtn = document.querySelector("#run-query")

function runQuery() {
    const editor = document.querySelector("#query-editor")

    runBtn.setAttribute("disabled", "")
    runBtn.classList.add("disabled")
    
    fetch(`/api/query?language=${language(editor)}`, {
        method: "POST",
        headers: {
            "Content-Type": "text/sql",
            "Accept": "application/json"
        },
        body: query(editor)
    })
    .then(res => res.json())
    .then(res => {
        const table = document.querySelector("#resultsTable")
        const err = document.querySelector("#resultError")
        if('error' in res) {
            table.classList.add("hidden")
            err.classList.remove("hidden")
            err.textContent = res.error
        } else {
            table.innerHTML = ""
            table.classList.remove("hidden")
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
        runBtn.removeAttribute("disabled")
        runBtn.classList.remove("disabled")
    })
}

addEventListener('DOMContentLoaded', () => {
    runBtn.addEventListener("click", runQuery)
})