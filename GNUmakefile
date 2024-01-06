default: testacc

# Generate code from json specification
.PHONY: framework-generator
framework-generator:
	tfplugingen-framework generate provider --input ./internal/provider/provider.json --output ./internal/provider --package provider

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
