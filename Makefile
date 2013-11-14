PROJECT_ROOT := $(shell pwd)
PROJECT_NAME := $(shell pwd | sed 's/.*\///g')
VENDOR_PATH := $(PROJECT_ROOT)/vendor
GOPATH := $(PROJECT_ROOT):$(VENDOR_PATH)

all: test fmt

fmt:
	@find . -name \*.go -exec gofmt -l -w {} \;

clean:
	@rm -rf bin pkg $(PROJECT_ROOT)/src/atlantis/crypto/key.go

copy-key:
	@cp $(ATLANTIS_SECRET_DIR)/atlantis_key.go $(PROJECT_ROOT)/src/atlantis/crypto/key.go

install-deps:
	@echo "Installing Dependencies..."
	@rm -rf $(VENDOR_PATH)
	@mkdir -p $(VENDOR_PATH) || exit 2
	@GOPATH=$(VENDOR_PATH) go get launchpad.net/gocheck
	@echo "Dependencies Installed."

test: clean copy-key
	@GOPATH=$(GOPATH) go test atlantis/common -gocheck.vv=true -test.v=true
	@GOPATH=$(GOPATH) go test atlantis/crypto -gocheck.vv=true -test.v=true
