install_deps:
	go get ./...
	go get github.com/cespare/reflex

start:
	go run cmd/main.go

dev:
	reflex --start-service -r '\.go$$' make start

go-mod-tidy: ## Check if go.mod and go.sum does not contains any unnecessary dependencies and remove them.
	$(call before_job,Go mod tidy checking dependencies:)
ifndef TMPDIR
	$(eval TMPDIR=$(shell mktemp -d))
endif
	cp -fv go.mod $(TMPDIR)
	cp -fv go.sum $(TMPDIR)
	go mod tidy -v
	diff -u $(TMPDIR)/go.mod go.mod
	diff -u $(TMPDIR)/go.sum go.sum
	rm -f $(TMPDIR)go.mod $(TMPDIR)go.sum
	$(call after_job,Go mod check succeeded!)
