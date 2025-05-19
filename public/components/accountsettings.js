addEventListener("DOMContentLoaded", () => {
    const loading = document.querySelector("#loading-page")

    document.querySelector("#refresh-db").addEventListener("click", async () => {
        loading.style.pointerEvents = "all";
        loading.style.opacity = 1;

        await fetch("/api/query?language=sql", {method: "post", body: "drop table pacuare_raw;"})
        await fetch("/api/refresh", {method: "post"})

        loading.style.pointerEvents = "none";
        loading.style.opacity = 0;
    })

    document.querySelector("#recreate-db").addEventListener("click", async () => {
        loading.style.pointerEvents = "all";
        loading.style.opacity = 1;

        await fetch("/api/recreate", {method: "post"})
        
        location.reload()
    })
})