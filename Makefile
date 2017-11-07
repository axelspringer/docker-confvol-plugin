PACKAGES=$(shell go list ./... | grep -v /vendor/)
RACE := $(shell test $$(go env GOARCH) != "amd64" || (echo "-race"))
VERSION := $(shell grep "const Version " version/const.go | sed -E 's/.*"(.+)"$$/\1/')
SRC = main.go
BINARY ?= docker-confvol-plugin
BINDIR ?= ./bin/

PLUGIN_INSTALL_DIR ?= $(DESTDIR)/var/lib/docker/
PLUGIN_DESC_INSTALL_DIR ?= $(DESTDIR)/etc/docker/
SYSTEM_INSTALL_DIR ?= $(DESTDIR)/usr/lib/systemd/system/
MAN_INSTALL_DIR ?= $(DESTDIR)/usr/share/man/

all: build man

help:
	@echo 'Available commands:'
	@echo
	@echo 'Usage:'
	@echo '    make deps     		Install go deps.'
	@echo '    make build    		Compile the project.'
	@echo '    make test    		Run ginkgo test suites.'
	@echo '    make man     		Create man doc'
	@echo '    make restore  		Restore all dependencies.'
	@echo '    make clean    		Clean the directory tree.'
	@echo

test: 
	ginkgo --cover -v driver
	go tool cover -html=driver/driver.coverprofile -o driver_test_coverage.html

deps:
	go get -u github.com/mitchellh/gox
	go get -u github.com/onsi/gomega
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/sirupsen/logrus
	go get -u github.com/docker/go-plugins-helpers/volume
	go get -u github.com/docker/libkv
	go get -u github.com/cpuguy83/go-md2man
	go get -u github.com/coreos/etcd/client

build: $(SRC)
	@echo "Compiling..."
	@mkdir -p "$(BINDIR)"
	go build -o "$(BINDIR)$(BINARY)" $^
	@echo "All done! The binaries is in ./bin let's have fun!"

vet: ## run go vet
	test -z "$$(go vet ${PACKAGES} 2>&1 | grep -v '*composite literal uses unkeyed fields|exit status 0)' | tee /dev/stderr)"

ci: vet test

man:
	go-md2man -in "files/man/$(BINARY).8.md" -out "files/man/$(BINARY).8"

install:
	install -D -m 644 "files/etc/docker/$(BINARY)" "$(PLUGIN_DESC_INSTALL_DIR)$(BINARY)"
	install -D -m 644 "files/systemd/$(BINARY).service" "$(SYSTEM_INSTALL_DIR)$(BINARY).service"
	install -D -m 644 "files/systemd/$(BINARY).socket" "$(SYSTEM_INSTALL_DIR)$(BINARY).socket"
	install -D -m 755 "$(BINDIR)$(BINARY)" "$(PLUGIN_INSTALL_DIR)$(BINARY)"
	install -D -m 644 "files/man/$(BINARY).8" "$(MAN_INSTALL_DIR)/man8/$(BINARY).8"

restore:
	dep ensure

.PHONY: install deps man restore