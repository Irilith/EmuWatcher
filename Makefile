OUTPUT := EmuWatcher.exe
SYSO := cmd/EmuWatcher/EmuWatcher.syso
ICON := assets/icon.ico
MANIFEST := assets/EmuWatcher.exe.xml
MAIN := main.go

UPDATER_DIR := updater
UPDATER_OUTPUT := Updater.exe
SOURCE_FILE := $(UPDATER_DIR)/target/release/$(UPDATER_OUTPUT)
DESTINATION_FILE := ./$(UPDATER_OUTPUT)

GO := go
CARGO := cargo
RSRC := rsrc

.PHONY: all build run build-updater move dev-build-updater dev-build clean

all: build build-updater move

$(SYSO): $(ICON) $(MANIFEST)
	$(RSRC) -ico $(ICON) -manifest $(MANIFEST) -o $(SYSO)

run: build
	./$(OUTPUT)

build: $(SYSO)
	@read -p "Enter version: " VERSION; \
	COMMIT_HASH=$$(git rev-parse --short HEAD); \
	$(GO) build -ldflags="-X EmuWatcher/utils/version.version=v$$VERSION -X EmuWatcher/utils/version.commit=$$COMMIT_HASH" -o $(OUTPUT) $(MAIN)

build-updater:
	@cd $(UPDATER_DIR) && $(CARGO) build --release

move:
	@if [ -f "$(SOURCE_FILE)" ]; then \
		rm -f "$(DESTINATION_FILE)"; \
		mv "$(SOURCE_FILE)" "$(DESTINATION_FILE)"; \
	else \
		echo "$(UPDATER_OUTPUT) not found!"; \
	fi

dev-build-updater:
	@cd $(UPDATER_DIR) && $(CARGO) build

dev-build: $(SYSO)
	$(GO) build -o $(OUTPUT) $(MAIN)

clean:
	rm -f $(OUTPUT) $(SYSO) $(DESTINATION_FILE)
	@cd $(UPDATER_DIR) && $(CARGO) clean
