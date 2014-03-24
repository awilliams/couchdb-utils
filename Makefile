BIN=couchdb-utils
VERSION=0.0.2
README=README.md
LICENSE=LICENSE
RELEASE_DIR=release

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GODEP=$(GOCMD) get -d -v ./... 
GOFMT=gofmt -w
 
default: build

build:
	$(GODEP)
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o bin/linux-amd64/$(BIN)
	GOARCH=386 GOOS=linux $(GOBUILD) -o bin/linux-386/$(BIN)
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o bin/darwin-amd64/$(BIN)

package:
	rm -rf $(RELEASE_DIR)/couchdb-utils
	mkdir $(RELEASE_DIR)/couchdb-utils 
	cp $(README) $(RELEASE_DIR)/couchdb-utils/$(README)
	cp $(LICENSE) $(RELEASE_DIR)/couchdb-utils/$(LICENSE)

	cp -f bin/linux-amd64/$(BIN) $(RELEASE_DIR)/couchdb-utils/$(BIN)
	tar -czf $(RELEASE_DIR)/$(BIN)-linux-amd64-v$(VERSION).tar.gz -C $(RELEASE_DIR) couchdb-utils

	cp -f bin/linux-386/$(BIN) $(RELEASE_DIR)/couchdb-utils/$(BIN)
	tar -czf $(RELEASE_DIR)/$(BIN)-linux-386-v$(VERSION).tar.gz -C $(RELEASE_DIR) couchdb-utils

	cp -f bin/darwin-amd64/$(BIN) $(RELEASE_DIR)/couchdb-utils/$(BIN)
	tar -czf $(RELEASE_DIR)/$(BIN)-darwin-amd64-v$(VERSION).tar.gz -C $(RELEASE_DIR) couchdb-utils

	rm -rf $(RELEASE_DIR)/couchdb-utils

format:
	$(GOFMT) ./**/*.go

clean:
	$(GOCLEAN)

test:
	$(GODEP) && $(GOTEST) -v ./...
