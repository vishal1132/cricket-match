run: build
	@./cmd/server/main

build:
	@go build -o ./cmd/server/main ./cmd/server

buildwasm:
	@GOOS=js GOARCH=wasm go build -o ./assets/js/match.wasm ./cmd/wasm