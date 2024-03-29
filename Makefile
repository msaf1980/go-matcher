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
	# glob
	${GO} test -coverprofile cover_glob.out -coverpkg=./glob,./pkg/items,./pkg/escape,./pkg/utils ./glob
	${GO} tool cover -html=cover_glob.out -o cover_glob.html
	# gglob
	${GO} test -coverprofile cover_gglob.out -coverpkg=./gglob,./glob,./pkg/items,./pkg/escape,./pkg/utils ./gglob
	${GO} tool cover -html=cover_gglob.out -o cover_gglob.html
	# gtags
	${GO} test -coverprofile cover_gtags.out -coverpkg=./gtags,./pkg/items,./pkg/escape ./gtags
	${GO} tool cover -html=cover_gtags.out -o cover_gtags.html
	# expand
	${GO} test -coverprofile cover_expand.out -coverpkg=./expand,./pkg/items,./pkg/escape ./expand
	${GO} tool cover -html=cover_expand.out -o cover_expand.html
	# ALL
	# gocovmerge cover_items.out cover_escape.out cover_utils.out cover_gglob.out cover_gtags.out > cover.out
	#${GO} test -coverprofile cover.out -coverpkg=./glob,./gglob,./gtags,./expand,./pkg/items,./pkg/escape,./pkg/utils ./...
	${GO} test -coverprofile cover.out -coverpkg=./gglob,./glob,./expand,./pkg/items,./pkg/escape,./pkg/utils ./...
	${GO} tool cover -html=cover.out -o cover.html

lint:
	golangci-lint run

FORCE:

.PHONY: build
