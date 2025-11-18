.PHONY: run-dataflow run-cfg run-mem run-ptr help all

help:
	@echo "Available targets:"
	@echo "  run-dataflow - Run dataflow analysis demo"
	@echo "  run-cfg      - Run CFG construction demo"
	@echo "  run-mem      - Run memory management demo"
	@echo "  run-ptr      - Run pointer and GC demo"
	@echo "  all          - Run all demos"

run-dataflow:
	go run ./cmd/tool dataflow

run-cfg:
	go run ./cmd/tool cfg

run-mem:
	go run ./cmd/tool mem

run-ptr:
	go run ./cmd/tool ptr

all: run-dataflow run-cfg run-mem run-ptr

