install_deps:
	go get ./...
	go get github.com/cespare/reflex

start:
	go run cmd/main.go

dev:
	reflex --start-service -r '\.go$$' make start

