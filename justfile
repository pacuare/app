templ:
    templ generate

build: templ
    go build -o out

dev:
    templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."
