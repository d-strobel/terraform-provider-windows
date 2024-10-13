# Terminal colors
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# Source .env file if available
ifneq ("$(wildcard .env)","")
	include .env
else
	include ./vagrant/vagrant.env
endif

# Generate all
.PHONY: generate
generate:
	@printf "$(OK_COLOR)==> Go generate$(NO_COLOR)\n"
	@go generate ./...

# Setup requirements
.PHONY: vagrant-up
vagrant-up:
	@printf "$(OK_COLOR)==> Setup vagrant machines$(NO_COLOR)\n"
	@$(MAKE) -C vagrant vagrant-up

# Remove requirements
.PHONY: vagrant-down
vagrant-down:
	@printf "$(OK_COLOR)==> Remove vagrant machines$(NO_COLOR)\n"
	@$(MAKE) -C vagrant vagrant-down

# Run acceptance tests
.PHONY: testacc
testacc: testacc-terraform testacc-opentofu

# Run Terraform acceptance tests
.PHONY: testacc-terraform
testacc-terraform:
	@printf "$(OK_COLOR)==> Run Terraform acceptance tests$(NO_COLOR)\n"
	TF_ACC=1 go test ./internal/provider/... -v $(TESTARGS) -timeout 120m

# Run OpenTofu acceptance tests
.PHONY: testacc-opentofu
testacc-opentofu:
	@printf "$(OK_COLOR)==> Run OpenTofu acceptance tests$(NO_COLOR)\n"
	TF_ACC_TERRAFORM_PATH="$(shell which tofu)" TF_ACC_PROVIDER_NAMESPACE="hashicorp" TF_ACC_PROVIDER_HOST="registry.opentofu.org" TF_ACC=1 go test ./internal/provider/... -v $(TESTARGS) -timeout 120m
