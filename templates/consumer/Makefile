.PHONY: help build check validate

help:
	@echo "  make build    - Build bin/validate from ai-sdlc/tools/validate"
	@echo "  make validate - Project gate: AC (all epics) + pipeline/structure per epic"
	@echo "  make check    - Product quality gate (go vet when go.mod exists)"

build: bin/validate

bin/validate:
	@test -f ai-sdlc/tools/validate/go.mod || (echo "Missing ai-sdlc/: clone process repo into ai-sdlc/ (see README)" >&2; exit 1)
	@mkdir -p bin
	go build -C ai-sdlc/tools/validate -o $(CURDIR)/bin/validate .

validate: bin/validate
	@if [ -z "$(filter-out validate,$(MAKECMDGOALS))" ]; then \
		./bin/validate; \
		set -e; for epic in ai-sdlc-artefacts/epics/EP-*; do \
			[ -d "$$epic" ] || continue; \
			id=$$(basename "$$epic"); \
			./bin/validate pipeline $$id; \
			./bin/validate structure $$id; \
		done; \
	else \
		./bin/validate $(filter-out validate,$(MAKECMDGOALS)); \
	fi

check:
	@if [ -f go.mod ]; then \
		go vet ./... && \
		BOOTSTRAP_MAKE_CHECK_INVOKED=1 go test ./tests/... -count=1; \
	else echo "check: no go.mod — skipped (add product module when ready)"; fi

%:
	@:
