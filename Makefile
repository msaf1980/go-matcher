all: test

DOCKER ?= docker
GO ?= go

## help: Prints a list of available build targets.
help:
	echo "Usage: make <OPTIONS> ... <TARGETS>"
	echo ""
	echo "Available targets are:"
	echo ''
	sed -n 's/^##//p' ${PWD}/Makefile | column -t -s ':' | sed -e 's/^/ /'
	echo
	echo "Targets run by default are: `sed -n 's/^all: //p' ./Makefile | sed -e 's/ /, /g' | sed -e 's/\(.*\), /\1, and /'`"

clean:
	rm cover*.out cover*.html

# prep:
# 	 ${GO} install github.com/wadey/gocovmerge

## test: Executes any unit tests.
test:
	${GO} test -cover -race ./...

# deep coverage (require Golang 1.20 or later)
coverage: FORCE
	${GO} test -coverprofile cover_gglob.out -coverpkg=./gglob,./pkg/items,./pkg/globs ./gglob
	${GO} tool cover -html=cover_gglob.out -o cover_gglob.html
	${GO} test -coverprofile cover_gtags.out -coverpkg=./gtags,./pkg/items,./pkg/escape ./gtags
	${GO} tool cover -html=cover_gtags.out -o cover_gtags.html
	# gocovmerge cover_items.out cover_escape.out cover_utils.out cover_gglob.out cover_gtags.out > cover.out
	${GO} test -coverprofile cover.out -coverpkg=./gglob,./gtags,./pkg/items,./pkg/globs,./pkg/escape ./...
	${GO} tool cover -html=cover.out -o cover.html

lint:
	golangci-lint run

FORCE:

.PHONY: build
