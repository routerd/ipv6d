# Copyright 2020 The routerd Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL=/bin/bash
.SHELLFLAGS=-euo pipefail -c

export CGO_ENABLED:=0

BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
SHORT_SHA=$(shell git rev-parse --short HEAD)
VERSION?=${BRANCH}-${SHORT_SHA}

# -------
# Compile
# -------
all: \
	all-mipsle \
	all-mips \
	all-mips64

all-mipsle: \
	bin/linux_mipsle/ipv6d

all-mips: \
	bin/linux_mips/ipv6d

all-mips64: \
	bin/linux_mips64/ipv6d

bin/linux_amd64/%: GOARGS = GOOS=linux GOARCH=amd64
bin/linux_mipsle/%: GOARGS = GOOS=linux GOARCH=mipsle
bin/linux_mips/%: GOARGS = GOOS=linux GOARCH=mips
bin/linux_mips64/%: GOARGS = GOOS=linux GOARCH=mips64

bin/%: FORCE
	$(eval COMPONENT=$(shell basename $*))
	$(GOARGS) go build -o bin/$* cmd/$(COMPONENT)/main.go

FORCE:

# -------
# Compile
# -------

clean:
	rm -rf bin/$*
.PHONY: clean

# -------------------
# Testing and Linting
# -------------------

test:
	CGO_ENABLED=1 go test -race -v ./...
.PHONY: test

fmt:
	go fmt ./...
.PHONY: fmt

vet:
	go vet ./...
.PHONY: vet

tidy:
	go mod tidy
.PHONY: tidy

verify-boilerplate:
	@go run hack/boilerplate/boilerplate.go \
		-boilerplate-dir hack/boilerplate/ \
		-verbose
.PHONY: verify-boilerplate

pre-commit-install:
	@echo "installing pre-commit hooks using https://pre-commit.com/"
	@pre-commit install
.PHONY: pre-commit-install
