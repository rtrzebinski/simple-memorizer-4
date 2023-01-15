build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm

run:
	make build
	go run main.go
