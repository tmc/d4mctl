.PHONY: all
all: install docs

.PHONY: lint
lint: lint-deps
	staticcheck ./...

.PHONY: lint-deps
lint-deps:
	@command -v staticcheck >/dev/null || GO111MODULE=off go get honnef.co/go/tools/cmd/staticcheck

.PHONY: install
install:
	go install

.PHONY: docs
docs:
	@d4mctl -h 2>&1 > .help
	@cat README.in | sh -c "HELP='$$(cat .help)' tmpl" > README.md

.PHONY: release
release:
	goreleaser --rm-dist
