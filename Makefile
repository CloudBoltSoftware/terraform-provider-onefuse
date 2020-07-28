PLUGIN_EXECUTABLE=terraform-provider-onefuse
TF_PLUGINS_DIR=$(HOME)/.terraform.d/plugins/

.PHONY : install
install: build
	mkdir -p $(TF_PLUGINS_DIR)
	mv -f $(PLUGIN_EXECUTABLE) $(TF_PLUGINS_DIR)

.PHONY : build
build:
	go build -o $(PLUGIN_EXECUTABLE)

.PHONY : clean
clean:
	go clean
	rm $(TF_PLUGINS_DIR)/$(PLUGINS_EXECUTABLE)
