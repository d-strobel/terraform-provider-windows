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

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./internal/provider/... -v $(TESTARGS) -timeout 120m
