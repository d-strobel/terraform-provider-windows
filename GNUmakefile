# Terminal colors
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# Generate code from json specification
.PHONY: framework-generator
framework-generator:
	@printf "$(OK_COLOR)==> Generate provider schema$(NO_COLOR)\n"
	@tfplugingen-framework generate provider --input ./internal/generator/provider_windows/provider_windows.json --output ./internal/generator/provider_windows --package provider_windows

	@printf "$(OK_COLOR)==> Generate local schema$(NO_COLOR)\n"
	@tfplugingen-framework generate data-sources --input ./internal/generator/local_datasources/local_datasources.json --output ./internal/generator/local_datasources --package local_datasources
	@tfplugingen-framework generate resources --input ./internal/generator/local_resources/local_resources.json --output ./internal/generator/local_resources --package local_resources

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
