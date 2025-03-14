.PHONY: templates

build:
	@go build -o ./bin/webrtc ./cmd/

run: build
	@echo 'Server: https://[::1]:8000'
	@./bin/webrtc

run-dev:
	@echo 'Local Server: https://localhost:8000'
	@go build -ldflags='-s -w' -o ./bin/webrtc ./cmd/
	@./bin/webrtc --addr=localhost:8000

watch:
	@templ generate -watch -proxy=http://localhost:8000
