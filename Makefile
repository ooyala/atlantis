## Copyright 2014 Ooyala, Inc. All rights reserved.
##
## This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
## except in compliance with the License. You may obtain a copy of the License at
## http://www.apache.org/licenses/LICENSE-2.0
##
## Unless required by applicable law or agreed to in writing, software distributed under the License is
## distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
## See the License for the specific language governing permissions and limitations under the License.

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
