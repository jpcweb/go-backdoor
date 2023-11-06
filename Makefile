BUILDER_BIN="backdoor"
build:
	@echo "Building binary..."
	env CGO_ENABLED=0 go build -o ./bin/${BUILDER_BIN} ./...
	@echo "Done!"
.PHONY: build
run: build
	@./bin/$(BUILDER_BIN)