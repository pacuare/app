import { apiQuery } from "../js/api.js";
import Alpine from "alpinejs";

document.addEventListener("alpine:init", () => {
    Alpine.data('accountSettings', () => ({
        loading: false,

        refreshButton: {
            async ["@click"]() {
                this.loading = true;

                await apiQuery("sql", "drop table pacuare_raw;");
                await fetch("/api/refresh", {method: "post"})

                this.loading = false;
            }
        },

        recreateButton: {
            async ["@click"]() {
                this.loading = true;

                await fetch("/api/recreate", {method: "post"})

                location.reload()
            }
        }
    }))
})