BIN=couchdb-utils
VERSION="0.0.1"
README=README.md
LICENSE=LICENSE
RELEASE_DIR=release

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOTEST) -i
GOFMT=gofmt -w
 
default: build

build:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o bin/linux-amd64/$(BIN)
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o bin/darwin-amd64/$(BIN)

package:
	tar -czf $(RELEASE_DIR)/$(BIN)-linux-amd64-v$(VERSION).tar.gz $(README) $(LICENSE) -C bin/linux-amd64 $(BIN)
	tar -czf $(RELEASE_DIR)/$(BIN)-darwin-amd64-v$(VERSION).tar.gz $(README) $(LICENSE) -C bin/darwin-amd64 $(BIN)

format:
	$(GOFMT) ./**/*.go

clean:
	$(GOCLEAN)
