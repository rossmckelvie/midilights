GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=ml

build:
	$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME) -v
