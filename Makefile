.PHONY: build

build:
	mkdir -p dist
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/plain-linux-amd64 ./cmd/plain
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/plain-windows-amd64.exe ./cmd/plain
	python3 support/plain-vscode/build_vsix.py
