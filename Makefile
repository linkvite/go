.PHONY: test lint fmt build clean cover generate check-generated release-major release-minor release-patch tag check-clean check-major proxy

# Get current version from version.go
VERSION := $(shell grep 'const Version' version.go | sed 's/.*"\(.*\)"/\1/')
MAJOR := $(shell echo $(VERSION) | cut -d. -f1)
MINOR := $(shell echo $(VERSION) | cut -d. -f2)
PATCH := $(shell echo $(VERSION) | cut -d. -f3)
# Major version the module path allows: /vN suffix in go.mod, else v0/v1
MODULE := $(shell sed -n 's/^module //p' go.mod)
MODULE_BASE := $(shell echo $(MODULE) | sed 's|/v[0-9][0-9]*$$||')
MODULE_MAJOR := $(or $(shell echo $(MODULE) | sed -n 's|.*/v\([0-9][0-9]*\)$$|\1|p'),1)

# Default target
all: lint test

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race ./...

# Run tests with coverage
cover:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Regenerate types_gen.go and services_gen.go from spec/openapi.json.
# SPEC=../cultus/api/openapi.json refreshes the spec first.
generate:
	@go run ./internal/gen $(if $(SPEC),-spec $(SPEC))

# Fail if the generated files don't match spec/openapi.json
check-generated:
	@go run ./internal/gen -check

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Apply gofumpt + goimports (the formatters lint checks)
fmt:
	@golangci-lint fmt

# Build (verify compilation)
build:
	@echo "Building..."
	@go build ./...

# Clean generated files
clean:
	@rm -f coverage.out coverage.html

# Show current version
version:
	@echo "Current version: $(VERSION)"

# Check if working tree and index are clean
check-clean:
	@if ! git diff --quiet; then \
		echo "Unstaged changes present. Commit/stash before releasing."; \
		git --no-pager diff --name-only; \
		exit 1; \
	fi
	@if ! git diff --cached --quiet; then \
		echo "Staged changes present. Commit/unstage before releasing."; \
		git --no-pager diff --cached --name-only; \
		exit 1; \
	fi
	@if [ -n "$$(git ls-files --others --exclude-standard)" ]; then \
		echo "Untracked files present. Add/ignore before releasing:"; \
		git ls-files --others --exclude-standard; \
		exit 1; \
	fi


# Release helpers
# A new major version needs a new module path (/v3, ...), so it isn't a
# bump-and-tag: Go rejects a v3.0.0 tag on a /v2 module. See RELEASE.md.
release-major:
	@echo "A major release changes the module path: $(MODULE) -> $(MODULE_BASE)/v$$(($(MODULE_MAJOR)+1))."
	@echo "Follow 'Major versions' in RELEASE.md (go.mod, imports, version.go), merge, then 'make tag'."
	@exit 1

release-minor: check-clean check-major
	@$(eval NEW_VERSION := $(MAJOR).$(shell echo $$(($(MINOR)+1))).0)
	@echo "Bumping version from $(VERSION) to $(NEW_VERSION)"
	@perl -pi -e 's/const Version = "$(VERSION)"/const Version = "$(NEW_VERSION)"/' version.go
	@git add version.go
	@git commit -m "chore: bump version to $(NEW_VERSION)"
	@git tag -a v$(NEW_VERSION) -m "Release v$(NEW_VERSION)"
	@echo "Created tag v$(NEW_VERSION)"
	@echo "Run 'git push && git push --tags' to publish"

release-patch: check-clean check-major
	@$(eval NEW_VERSION := $(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1))))
	@echo "Bumping version from $(VERSION) to $(NEW_VERSION)"
	@perl -pi -e 's/const Version = "$(VERSION)"/const Version = "$(NEW_VERSION)"/' version.go
	@git add version.go
	@git commit -m "chore: bump version to $(NEW_VERSION)"
	@git tag -a v$(NEW_VERSION) -m "Release v$(NEW_VERSION)"
	@echo "Created tag v$(NEW_VERSION)"
	@echo "Run 'git push && git push --tags' to publish"

# Tag the version already in version.go, without bumping it (e.g. the first
# release of a new major: version.go says 2.0.0, nothing to bump)
tag: check-clean check-major
	@if git rev-parse -q --verify "refs/tags/v$(VERSION)" >/dev/null; then echo "v$(VERSION) already exists"; exit 1; fi
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@echo "Created tag v$(VERSION) on $$(git rev-parse --short HEAD)"
	@echo "Run 'git push origin v$(VERSION)' to publish"

# version.go's major must match the module path, or Go rejects the tag
check-major:
	@if [ "$(MAJOR)" != "$(MODULE_MAJOR)" ] && ! { [ "$(MODULE_MAJOR)" = 1 ] && [ "$(MAJOR)" = 0 ]; }; then \
		echo "version.go says $(VERSION) but $(MODULE) only takes v$(MODULE_MAJOR).x.y tags."; \
		echo "Fix version.go, or for a new major follow 'Major versions' in RELEASE.md."; \
		exit 1; \
	fi

# Trigger the Go module proxy to index the current version
proxy: check-major
	@echo "Triggering Go module proxy for v$(VERSION)..."
	@GOPROXY=proxy.golang.org go list -m $(MODULE)@v$(VERSION)
	@echo "Done. v$(VERSION) is now available via go get."

# Pre-release checks
check: check-generated lint test
	@echo "All checks passed!"

# Help
help:
	@echo "Available targets:"
	@echo "  make generate      - Regenerate from spec/openapi.json (SPEC=path to refresh it)"
	@echo "  make check-generated - Fail if generated code is stale"
	@echo "  make test          - Run tests"
	@echo "  make cover         - Run tests with coverage report"
	@echo "  make lint          - Run golangci-lint"
	@echo "  make fmt           - Format the code (gofumpt, goimports)"
	@echo "  make build         - Verify compilation"
	@echo "  make clean         - Remove generated files"
	@echo "  make version       - Show current version"
	@echo "  make check         - Run lint and test"
	@echo "  make check-clean   - Ensure working tree is clean before releasing"
	@echo ""
	@echo "Release targets:"
	@echo "  make release-major - Explains a major release (needs a new module path)"
	@echo "  make release-minor - Bump minor version (0.x.0)"
	@echo "  make release-patch - Bump patch version (0.0.x)"
	@echo "  make tag           - Tag the current version.go version without bumping"
	@echo "  make proxy         - Trigger Go module proxy to index current version"
