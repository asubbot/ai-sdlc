.PHONY: help build check validate

help:
	@echo "  make build    - Build bin/validate from ai-sdlc/tools/validate"
	@echo "  make validate - Run validator (e.g. make validate EP-000, make validate pipeline EP-000)"
	@echo "  make check    - Product quality gate (go vet when go.mod exists)"

build: bin/validate

bin/validate:
	@test -f ai-sdlc/tools/validate/go.mod || (echo "Missing ai-sdlc/: clone process repo into ai-sdlc/ (see README)" >&2; exit 1)
	@mkdir -p bin
	go build -C ai-sdlc/tools/validate -o $(CURDIR)/bin/validate .

validate: bin/validate
	@./bin/validate $(filter-out validate,$(MAKECMDGOALS))

check:
	@if [ -f go.mod ]; then go vet ./...; else echo "check: no go.mod — skipped (add product module when ready)"; fi

%:
	@:
