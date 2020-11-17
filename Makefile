.PHONY: all
all: install docs

.PHONY: install
install:
	go install

.PHONY: docs
docs:
	d4mctl -h 2>&1 > .help
	cat README.in | sh -c "HELP='$$(cat .help)' tmpl" > README.md 

.PHONY: release
release:
	goreleaser
