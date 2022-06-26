DIR           =$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

GOPATH       := $(DIR):$(GOPATH)
DATE          =$(shell date -u +%Y%m%d.%H%M%S.%Z)
TESTPACKETS   =$(shell if [ -f .testpackages ]; then cat .testpackages; fi)
BENCHPACKETS  =$(shell if [ -f .benchpackages ]; then cat .benchpackages; fi)

default: lint test

link:
	@echo "prepare..."
	@mkdir src 2>/dev/null; true
	@if [ ! -L $(DIR)/src/job ]; then ln -s $(DIR)/job $(DIR)/src/job 2>/dev/null; fi
	@if [ ! -L $(DIR)/src/pool ]; then ln -s $(DIR)/pool $(DIR)/src/pool 2>/dev/null; fi
	@if [ ! -L $(DIR)/src/types ]; then ln -s $(DIR)/types $(DIR)/src/types 2>/dev/null; fi
	@if [ ! -L $(DIR)/src/event ]; then ln -s $(DIR)/event $(DIR)/src/event 2>/dev/null; fi
	@cd ${DIR}/src && ln -s . github.com 2>/dev/null; true
	@cd ${DIR}/src && ln -s . webnice 2>/dev/null; true
	@cd ${DIR}/src && ln -s . job 2>/dev/null; true
	@if command -v "gvt" >/dev/null; then cd ${DIR}/src; GOPATH="$(DIR)" gvt fetch -branch v2 "github.com/webnice/lv2" 2>/dev/null; true; fi
	@if command -v "gvt" >/dev/null; then cd ${DIR}/src; GOPATH="$(DIR)" gvt fetch -branch v1 "github.com/webnice/debug" 2>/dev/null; true; fi
.PHONY: link

test: link
	@echo "mode: set" > $(DIR)/coverage.log
	@for PACKET in $(TESTPACKETS); do \
		touch $(DIR)/coverage-tmp.log; \
		unset GOPATH; go test -v -covermode=count -coverprofile=$(DIR)/coverage-tmp.log $$PACKET; \
		if [ "$$?" -ne "0" ]; then exit $$?; fi; \
		tail -n +2 $(DIR)/coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> $(DIR)/coverage.log; \
		rm -f $(DIR)/coverage-tmp.log; true; \
	done
.PHONY: test

cover: test
	unset GOPATH; go tool cover -html=$(DIR)/coverage.log
	@make clean
.PHONY: cover

bench: link
	@for PACKET in $(BENCHPACKETS); do GOPATH=${GOPATH} go test -race -bench=. -benchmem $$PACKET; done
	@make clean
.PHONY: bench

lint: link
	GOPATH=${GOPATH} gometalinter \
	--vendor \
	--deadline=15m \
	--cyclo-over=15 \
	--disable=aligncheck \
	--linter="vet:go tool vet -printf {path}/*.go:PATH:LINE:MESSAGE" \
	./...
	@make clean
.PHONY: lint

clean:
	@echo "cleaning..."
	@rm -rf ${DIR}/src; true
	@rm -rf ${DIR}/bin/*; true
	@rm -rf ${DIR}/pkg/*; true
	@rm -rf ${DIR}/*.log; true
.PHONY: clean
