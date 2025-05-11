MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:

.PHONY: audit
audit: tidy fmt
	go vet ./...
	go tool -modfile=go.tool.mod staticcheck ./...
	go tool -modfile=go.tool.mod govulncheck ./...
	golangci-lint run -v


.PHONY: build
build:

	goreleaser build --clean --single-target --snapshot

.PHONY: clean
clean:
ifneq (,$(wildcard ./dist))
	rm -rf dist/

endif

ifneq (,$(wildcard ./coverage))
	rm -rf coverage/

endif

.PHONY: container
container: tidy
	./scripts/build-container.sh

.PHONY: fmt
fmt:
	go tool -modfile=go.tool.mod golines --base-formatter=gofumpt -w .
	go tool -modfile=go.tool.mod gofumpt -l -w -extra .

.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: test
test:
	go test ./... -cover

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: update
update:
	go get -u ./...
	go mod tidy

# annoyed-aardvark start
.PHONY: install
install: build

	cp dist/gh-tp_linux_amd64_v1/annoyed-aardvark .

	gh ext remove aa

	gh ext install .

	-gh aa --version

# annoyed-aardvark stop

# annoyed-aardvark start
.PHONY: moreclean
moreclean: clean

ifneq (,$(wildcard ./.tp.toml*))
	rm .tp.toml*
endif

ifneq (,$(wildcard ~/.tp.toml*))
	rm ~/.tp.toml*
endif

ifneq (,$(wildcard ~/.config/*tp.toml*))
	rm ~/.config/*tp.toml*
endif

ifneq (,$(wildcard /var/tmp/*tp.toml*))
	rm /var/tmp/*tp.toml*
endif

ifneq (,$(wildcard ~/.config/gh-tp/.tp.toml*))
	rm -rf ~/.config/gh-tp
endif

.PHONY: test
test: container
	@mkdir -p coverdata
	@docker run --rm -v $(PWD)/coverdata:/app/coverdata esacteksab/tpt:$(shell cat .current-tag)
	@if [ -f "coverdata/coverage.out" ]; then \
		go tool cover -html=coverdata/coverage.out -o coverdata/coverage.html; \
		echo "Coverage report generated: coverage.html"; \
	else \
		echo "Error: Coverage file not found"; \
		exit 1; \
	fi

# annoyed-aardvark stop

# annoyed-aardvark start
.PHONY: update
update:
	go get -u ./...
	go get -u -modfile=go.tool.mod tool
	go mod tidy
# annoyed-aardvark stop