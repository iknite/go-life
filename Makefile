.DEFAULT_GOAL := install

# AutoEnv
ifeq ($(CI),)
	ENV ?= .env
	ENV_GEN := $(shell ./.env.gen ${ENV} .env.required)
	include ${ENV}
	export $(shell sed 's/=.*//' ${ENV}))
endif

.PHONY: test
test:
	@go test -v -count=1 ./...


.PHONY: lint
lint:
	@golangci-lint run


.PHONY: prepare
prepare:
	@rm -rf dist
	@go mod download
	@go generate github.com/iknite/raft-life/preset


.PHONY: install
install: prepare
	@go install


release-%: prepare
	@bumpversion $*
	@git push --follow-tags
	@gorelease
