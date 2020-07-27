PKGNAME=onefuse
PLUGIN_EXECUTABLE=terraform-provider-$(PKGNAME)
VERSION=$$(cat VERSION)
PLUGIN_RELEASE_EXECUTABLE=$(PLUGIN_EXECUTABLE)_v$(VERSION)
TF_PLUGINS_DIR=$(HOME)/.terraform.d/plugins/

GOFMT_FILES?=$$(find . -name '*.go' -not -path './vendor/*' -not -path './go/*')

default: build

# Build the plugin
build: fmtcheck
	go install

# Build the provider and copy it to your local terraform plugins directory for local integratin testing
install: fmtcheck
	go build -o $(PLUGIN_EXECUTABLE)
	mv $(PLUGIN_EXECUTABLE) $(TF_PLUGINS_DIR)

# Format code
fmt:
	gofmt -w $(GOFMT_FILES)

# Verify code conforms to gofmt standards
fmtcheck:
	test -n $$(gofmt -l $(GOFMT_FILES))

release-%: fmtcheck
	@sh -c "'$(CURDIR)/scripts/build.sh' --$* --sha256sum --output $(PLUGIN_RELEASE_EXECUTABLE) --basedir release/terraform-provider-onefuse"

release: release-darwin release-linux release-windows

clean:
	@rm -rf release/*

.PHONY : build clean install fmt fmtcheck
