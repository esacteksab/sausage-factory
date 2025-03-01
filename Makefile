MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:

.PHONY: clean
clean:
ifneq (,$(wildcard ./plan.md))
	rm plan.md
endif

ifneq (,$(wildcard ./plan.out))
	rm plan.out
endif

ifneq (,$(wildcard ./terraform.tfstate))
	rm terraform.tfstate
endif

ifneq (,$(wildcard ./bsTF))
	rm -r bsTF
endif
