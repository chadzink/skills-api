.PHONY: up down up-pg up-airflow

test:
	go test ./tests/... -v

test-short:
	go test ./tests/... -v -test.short

build:
	go build -buildvcs=false -o ./tmp/main .