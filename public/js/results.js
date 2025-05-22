import { apiQuery } from "./api.js"
import Alpine from "alpinejs"
import { downloadResults } from "./index/download.js";

Alpine.store('results', {
    currentResults: [],
    runDisabled: false,

    async runQuery(sql) {
        this.runDisabled = true;
        this.currentResults = await apiQuery('sql', sql);
        this.runDisabled = false;
    },

    download() {
        downloadResults(this.currentResults);
    }
})