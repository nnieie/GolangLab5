DIR = $(shell pwd)
CMD = $(DIR)/cmd
CONFIG_PATH = $(DIR)/config
IDL_PATH = $(DIR)/idl
OUTPUT_PATH = $(DIR)/output
API_PATH= $(DIR)/cmd/api

MODULE = github.com/nnieie/golanglab5

SERVICES := api user


.PHONY: kitex-gen-%
kitex-gen-%:
	mkdir -p $(CMD)/$* && cd $(CMD)/$* && \
	kitex \
	-gen-path ../../kitex_gen \
	-service "$*" \
	-module "$(MODULE)" \
	-type thrift \
	$(DIR)/idl/$*.thrift
	go mod tidy


.PHONY: kitex-update-%
kitex-update-%:
	kitex -module "${MODULE}" idl/$*.thrift


.PHONY: hertz-gen-api
hertz-gen-api:
	cd ${API_PATH}; \
	hz update -idl ${IDL_PATH}/api.thrift; \
	go mod tidy

.PHONY: start-%
start-%:
	go run ./cmd/$* --log-level=debug