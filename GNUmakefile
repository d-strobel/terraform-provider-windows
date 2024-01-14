# Terminal colors
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# Generate code from json specification
.PHONY: generate-framework
generate-framework:
	@printf "$(OK_COLOR)==> Generate provider schema$(NO_COLOR)\n"
	tfplugingen-framework generate provider --input ./internal/schema/provider_windows.json --output ./internal/generate/provider_windows --package provider_windows

	@printf "$(OK_COLOR)==> Generate local schema$(NO_COLOR)\n"
	tfplugingen-framework generate data-sources --input ./internal/schema/local_datasources.json --output ./internal/generate/local_datasources --package local_datasources
	tfplugingen-framework generate resources --input ./internal/schema/local_resources.json --output ./internal/generate/local_resources --package local_resources

# Generate documentation
.PHONY: generate-docs
generate-docs:
	@printf "$(OK_COLOR)==> Generate documentation$(NO_COLOR)\n"
	go generate ./...

# Generate all
.PHONY: generate
generate: generate-framework generate-docs

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
