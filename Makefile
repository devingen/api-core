# Check if VERSION is provided
ifndef VERSION
$(error VERSION is not defined. Use 'make archive-release VERSION=<version>' to provide it.)
endif

# Build for multiple platforms with version
.PHONY: release
release:
	git tag -a v$(VERSION) -m ""
	git push origin v$(VERSION)

# Archive the binaries into .tar.gz and .zip files
.PHONY: archive-release
archive-release: clean-release build-release
	cp $(OUTPUT_DIR)/damga-linux-amd64 $(OUTPUT_DIR)/damga
	tar -czvf $(OUTPUT_DIR)/damga-cli-linux-amd64.tar.gz -C $(OUTPUT_DIR) damga
	rm $(OUTPUT_DIR)/damga $(OUTPUT_DIR)/damga-linux-amd64

	cp $(OUTPUT_DIR)/damga-windows-amd64.exe $(OUTPUT_DIR)/damga.exe
	zip -j $(OUTPUT_DIR)/damga-cli-windows-amd64.zip $(OUTPUT_DIR)/damga.exe
	rm $(OUTPUT_DIR)/damga.exe $(OUTPUT_DIR)/damga-windows-amd64.exe

	cp $(OUTPUT_DIR)/damga-darwin-amd64 $(OUTPUT_DIR)/damga
	tar -czvf $(OUTPUT_DIR)/damga-cli-darwin-amd64.tar.gz -C $(OUTPUT_DIR) damga
	rm $(OUTPUT_DIR)/damga $(OUTPUT_DIR)/damga-darwin-amd64

.PHONY: clean-release
clean-release:
	rm -rf $(OUTPUT_DIR)