GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=tunacan

build:
	$(GOBUILD) -o $(BINARY_NAME) cmd/tunacan/main.go
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
deps:
	$(GOGET) -u cloud.google.com/go/storage
	$(GOGET) github.com/mitchellh/cli
