ENGINE_DIR=cmd

BUILD_DIR=.bin

binaries:
	@echo "Building lambda binaries"
	rm -rf ./bin
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o ${BUILD_DIR}/book ${ENGINE_DIR}/main.go
	@echo "Success build app. Your app is ready to use in 'bin/' directory."
.PHONY: binaries

zip:
	zip -j .bin/book.zip .bin/book
.PHONY: zip


build: binaries zip
.PHONY: build

start:
	sls offline --useDocker start --host 0.0.0.0