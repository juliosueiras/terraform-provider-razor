all: build

install:
	go install

build:
	go build ; cp terraform-provider-razor ~/.terraform.d/plugins/terraform-provider-razor; rm -rf .terraform/; terraform init

.PHONY: install
