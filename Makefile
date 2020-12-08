# IMPORTANT: ⚠️ Requires gmake >= 4.0 ⚠️

PKGNAME := onefuse
PLUGIN_EXECUTABLE := terraform-provider-$(PKGNAME)
VERSION := $(file < VERSION)  # `file` may be a Make 4.3+ feature
ifeq ($(OS),Windows_NT)
	PLUGIN_RELEASE_EXECUTABLE := $(strip $(PLUGIN_EXECUTABLE)_v$(VERSION)).exe
else
	PLUGIN_RELEASE_EXECUTABLE := $(strip $(PLUGIN_EXECUTABLE)_v$(VERSION))
endif
VERSION_NUM := $(firstword $(subst -, ,$(VERSION)))
HOSTOS := $$(go env GOHOSTOS)
HOSTARCH := $$(go env GOHOSTARCH)
PLUGIN_RELEASE_EXECUTABLE := $(PLUGIN_EXECUTABLE)_v$(VERSION)
TF_PLUGINS_DIR_0.12 := $(HOME)/.terraform.d/plugins/$(HOSTOS)_$(HOSTARCH)# TODO: Drop support for TF 0.12
TF_PLUGINS_DIR := $(HOME)/.terraform.d/plugins/cloudbolt.io/terraform/$(PKGNAME)/$(VERSION_NUM)/$(HOSTOS)_$(HOSTARCH)

GOFMT_FILES?=$$(find . -name '*.go' -not -path './vendor/*' -not -path './go/*')

default: build

# Build the plugin
build: fmtcheck
	go install
	go build -o $(PLUGIN_EXECUTABLE)

install: build
	mkdir -p $(TF_PLUGINS_DIR)
	cp -f $(PLUGIN_EXECUTABLE) $(TF_PLUGINS_DIR)
	# Terraform 0.12 compatability
	mkdir -p $(TF_PLUGINS_DIR_0.12) # TODO: Drop support for TF 0.12
	cp -f $(PLUGIN_EXECUTABLE) $(TF_PLUGINS_DIR_0.12) # TODO: Drop support for TF 0.12

# Format code
fmt:
	gofmt -w main.go
	gofmt -w onefuse

# Verify code conforms to gofmt standards
fmtcheck:
	@gofmt -l main.go
	@gofmt -l onefuse
ifneq ($(strip $(gofmt -l main.go)),)
	@exit 1
endif
ifneq ($(strip $(gofmt -l onefuse)),)
	@exit 1
endif

release-%: fmtcheck
	scripts/build.sh --$* --sha256sum --output $(PLUGIN_RELEASE_EXECUTABLE) --basedir release/terraform-provider-onefuse

release: release-darwin release-linux release-windows

clean:
	@rm -rf release/*

# Live API tests
testacc:
	cd onefuse ; source config.env ; go test

.PHONY : build clean install fmt fmtcheck
