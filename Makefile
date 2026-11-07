.PHONY: test lint build clean cover generate check-generated release-major release-minor release-patch tag check-clean proxy

# Get current version from version.go
VERSION := $(shell grep 'const Version' version.go | sed 's/.*"\(.*\)"/\1/')
MAJOR := $(shell echo $(VERSION) | cut -d. -f1)
MINOR := $(shell echo $(VERSION) | cut -d. -f2)
PATCH := $(shell echo $(VERSION) | cut -d. -f3)

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
release-major: check-clean
	@$(eval NEW_VERSION := $(shell echo $$(($(MAJOR)+1)).0.0))
	@echo "Bumping version from $(VERSION) to $(NEW_VERSION)"
	@perl -pi -e 's/const Version = "$(VERSION)"/const Version = "$(NEW_VERSION)"/' version.go
	@git add version.go
	@git commit -m "chore: bump version to $(NEW_VERSION)"
	@git tag -a v$(NEW_VERSION) -m "Release v$(NEW_VERSION)"
	@echo "Created tag v$(NEW_VERSION)"
	@echo "Run 'git push && git push --tags' to publish"

release-minor: check-clean
	@$(eval NEW_VERSION := $(MAJOR).$(shell echo $$(($(MINOR)+1))).0)
	@echo "Bumping version from $(VERSION) to $(NEW_VERSION)"
	@perl -pi -e 's/const Version = "$(VERSION)"/const Version = "$(NEW_VERSION)"/' version.go
	@git add version.go
	@git commit -m "chore: bump version to $(NEW_VERSION)"
	@git tag -a v$(NEW_VERSION) -m "Release v$(NEW_VERSION)"
	@echo "Created tag v$(NEW_VERSION)"
	@echo "Run 'git push && git push --tags' to publish"

release-patch: check-clean
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
tag: check-clean
	@if git rev-parse -q --verify "refs/tags/v$(VERSION)" >/dev/null; then echo "v$(VERSION) already exists"; exit 1; fi
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@echo "Created tag v$(VERSION) on $$(git rev-parse --short HEAD)"
	@echo "Run 'git push origin v$(VERSION)' to publish"

# Trigger the Go module proxy to index the current version
proxy:
	@echo "Triggering Go module proxy for v$(VERSION)..."
	@GOPROXY=proxy.golang.org go list -m github.com/linkvite/go/v2@v$(VERSION)
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
	@echo "  make build         - Verify compilation"
	@echo "  make clean         - Remove generated files"
	@echo "  make version       - Show current version"
	@echo "  make check         - Run lint and test"
	@echo "  make check-clean   - Ensure working tree is clean before releasing"
	@echo ""
	@echo "Release targets:"
	@echo "  make release-major - Bump major version (x.0.0)"
	@echo "  make release-minor - Bump minor version (0.x.0)"
	@echo "  make release-patch - Bump patch version (0.0.x)"
	@echo "  make tag           - Tag the current version.go version without bumping"
	@echo "  make proxy         - Trigger Go module proxy to index current version"
