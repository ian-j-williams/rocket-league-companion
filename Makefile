.PHONY: test build build-linux build-windows build-mac clean

GO ?= go
BIN_DIR := dist

test:
	$(GO) test ./...

build:
	mkdir -p $(BIN_DIR)
	$(GO) build -trimpath -ldflags "-s -w" -o $(BIN_DIR)/rocketleague-companion ./cmd/companion

build-linux:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -trimpath -ldflags "-s -w" -o $(BIN_DIR)/rocketleague-companion-linux-amd64 ./cmd/companion

build-windows:
	mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags "-s -w" -o $(BIN_DIR)/rocketleague-companion-windows-amd64.exe ./cmd/companion

build-mac:
	mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags "-s -w" -o $(BIN_DIR)/rocketleague-companion-darwin-arm64 ./cmd/companion

clean:
	rm -rf $(BIN_DIR)
