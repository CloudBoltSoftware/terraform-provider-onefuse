PKGNAME := onefuse
PLUGIN_EXECUTABLE := terraform-provider-$(PKGNAME)
VERSION := $(strip $(file < VERSION))  # `file` may be a Make 4.3+ feature
ifeq ($(OS),Windows_NT)
	PLUGIN_RELEASE_EXECUTABLE := $(PLUGIN_EXECUTABLE)_v$(VERSION).exe
else
	PLUGIN_RELEASE_EXECUTABLE := $(PLUGIN_EXECUTABLE)_v$(VERSION)
endif
TF_PLUGINS_DIR := $$HOME/.terraform.d/plugins

default: build

# Build the plugin
install:
	go install

# Build the provider and copy it to your local terraform plugins directory for local integratin testing
build: install
	go build -o $(PLUGIN_RELEASE_EXECUTABLE)
	echo Move $(PLUGIN_RELEASE_EXECUTABLE) to $(TF_PLUGINS_DIR)

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

.PHONY : build clean install fmt fmtcheck
