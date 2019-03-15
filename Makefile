install_deps:
	go get ./...
	go get github.com/cespare/reflex

start:
	go run cmd/main.go

dev:
	reflex --start-service -r '\.go$$' make start

# Colorful output
color_off = \033[0m
color_cyan = \033[1;36m
color_green = \033[1;32m

define log_info
	@printf "$(color_cyan)$(1)$(color_off)\n"
endef
define log_success
	@printf "$(color_green)$(1)$(color_off)\n"
endef

test:
	$(call log_info,Run tests and check race conditions)
	# https://golang.org/doc/articles/race_detector.html
	go test -race -v ./...
	$(call log_success,All tests succeeded)

go-mod-tidy:
	$(call log_info,Check that go.mod and go.sum don't contain any unnecessary dependency)
	$(eval TMPDIR=$(shell mktemp -d))
	cp -f go.mod $(TMPDIR)
	cp -f go.sum $(TMPDIR)
	go mod tidy -v
	diff -u $(TMPDIR)/go.mod go.mod
	diff -u $(TMPDIR)/go.sum go.sum
	rm -rf $(TMPDIR)
	$(call log_success,Go mod check succeeded!)

test/ci:
	make test
	make go-mod-tidy

test/watch:
	reflex --start-service -r '\.go$$' make test

