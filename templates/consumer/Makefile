.PHONY: help fmt vet vuln lint test-race-cover coverage coverage-html \
	test test-race test-integration test-e2e check-boundaries \
	build validate check

GOLANGCI_LINT_VERSION ?= v2.5.0

help:
	@echo "  make build     - Build bin/validate from ai-sdlc/tools/validate"
	@echo "  make validate  - Product gate: AC + ears + req (all epics; ears/req skip NEW/CANCEL) + pipeline/structure per epic"
	@echo "  make check     - Product gate: fmt, vet, vuln, lint, race tests + coverage on tests/"
	@echo ""
	@echo "  Uncomment Makefile targets marked for EP-001+ when cmd/, integration, or e2e tests exist."

fmt:
	go fmt ./...

vet:
	go vet ./...

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run ./...

# Active on greenfield skeleton (tests/ only). BOOTSTRAP_MAKE_CHECK_INVOKED avoids recursive make check in bootstrap_test.go.
test-race-cover:
	@rm -f coverage.out
	BOOTSTRAP_MAKE_CHECK_INVOKED=1 go test -race ./tests/... -count=1 -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

check:
	@if [ ! -f go.mod ]; then echo "check: no go.mod — skipped" >&2; exit 0; fi
	@$(MAKE) fmt vet vuln lint test-race-cover

# --- Enable when product layout grows (uncomment targets and add to check:) ---

# test: Run all packages (replace test-race-cover scope). Enable when cmd/ exists.
# test:
# 	go test -tags=integration ./...

# test-race: Full module with race detector. Enable with test on ./...
# test-race:
# 	go test -race -tags=integration ./...

# test-integration: Requires Docker and //go:build integration tests.
# test-integration:
# 	go test -tags=integration ./tests/integration/...

# test-e2e: Requires tests/e2e/ and //go:build e2e (or integration,e2e).
# test-e2e:
# 	go test -tags=integration,e2e -count=1 ./tests/e2e/...

# coverage: Full-module coverage. Enable when ./... tests exist.
# coverage:
# 	rm -f coverage.out
# 	go test -tags=integration -count=1 ./... -coverpkg=./... -coverprofile=coverage.out -covermode=atomic
# 	go tool cover -func=coverage.out

# coverage-html: HTML report from full-module coverage.
# coverage-html: coverage
# 	go tool cover -html=coverage.out -o coverage.html

# check-boundaries: Module import rules. Enable when cmd/ and internal/ exist; see scripts/check-module-boundaries.sh.
# check-boundaries:
# 	./scripts/check-module-boundaries.sh

# git diff --exit-code: Enable in ci.yml when the repo has an initial commit (fails in fresh dirs after fmt).

build: bin/validate

bin/validate:
	@test -f ai-sdlc/tools/validate/go.mod || (echo "Missing ai-sdlc/: clone process repo into ai-sdlc/ (see README)" >&2; exit 1)
	@mkdir -p bin
	go build -C ai-sdlc/tools/validate -o $(CURDIR)/bin/validate .

validate: bin/validate
	@./bin/validate $(filter-out validate,$(MAKECMDGOALS))
	@if [ -z "$(filter-out validate,$(MAKECMDGOALS))" ]; then \
		set -e; for epic in ai-sdlc-artefacts/epics/EP-*; do \
			[ -d "$$epic" ] || continue; \
			id=$$(basename "$$epic"); \
			./bin/validate pipeline $$id; \
			./bin/validate structure $$id; \
		done; \
	fi

%:
	@:
