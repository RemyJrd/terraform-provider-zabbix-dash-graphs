# Variables
NAME := zabbix-dash-graphs
VERSION := 1.1.5
DIST := dist
PKG := ./zabbix
PLUGIN_DIR := ~/.terraform.d/plugins

OS_ARCHES := \
    linux_amd64 \
    darwin_amd64 \
    darwin_arm64 \
    windows_amd64

default: build

.PHONY: all build install uninstall test testacc clean release

build:
	mkdir -p $(DIST)
	@for target in $(OS_ARCHES); do \
		OS=$${target%_*}; ARCH=$${target#*_}; \
		BIN=terraform-provider-$(NAME)_v$(VERSION); \
		ZIP=terraform-provider-$(NAME)_$(VERSION)_$$OS\_$$ARCH.zip; \
		echo "Building for $$OS/$$ARCH..."; \
		CGO_ENABLED=0 GOOS=$$OS GOARCH=$$ARCH go build -o $(DIST)/$$BIN ./; \
		zip -j -q $(DIST)/$$ZIP $(DIST)/$$BIN; \
		rm $(DIST)/$$BIN; \
		echo "Built: $(DIST)/$$ZIP"; \
	done	


install:
	mkdir -p $(PLUGIN_DIR)
	go build -o $(PLUGIN_DIR)/terraform-provider-$(NAME)

uninstall:
	@rm -vf $(PLUGIN_DIR)/terraform-provider-$(NAME)

# Tests unitaires
test:
	go test $(PKG) || exit 1
	echo $(PKG) | xargs -t -n4 go test -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(PKG) -v -timeout 120m

clean:
	rm -rf $(DIST)

