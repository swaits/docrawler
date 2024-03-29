# Golang project Makefile
#
# derived from: http://zduck.com/2014/go-project-structure-and-dependencies/

# build targets (which aren't files)
.PHONY: build fmt lint vet test autotest cover doc run clean env sh zsh stats vendor

# configuration
APPNAME := docrawler
VENDOR := .vendor
PACKAGES :=                     \
	golang.org/x/tools/cmd/cover  \
	golang.org/x/tools/cmd/godoc  \
	golang.org/x/tools/cmd/vet    \
	github.com/golang/lint/golint

# override GOPATH to be our vendor directory
export GOPATH := ${PWD}/${VENDOR}:${PWD}
export PATH := ${PWD}/${VENDOR}/bin:${PWD}/bin:${PATH}

default: build

build: fmt lint vet
	go build -v -o ./bin/${APPNAME} ./src/${APPNAME}

fmt:
	go fmt ./src/...

# TODO golint is not accepting the "./src/..." for some reason (from within make)
lint:
	golint ./src/${APPNAME}/*.go

vet:
	go vet ./src/...

test:
	go test -v ./src/... | ./util/testfilter.sh

autotest:
	fswatch -or ./src | \
		xargs -n1 -I{} ${SHELL} -c 'echo ; date ; echo ----------------------------; go test ./src/... | ./util/testfilter.sh'

cover:
	mkdir -p tmp
	go test -coverprofile tmp/cover.out ./src/...
	go tool cover -html tmp/cover.out

doc:
	godoc -http=:6060 -index

run: build
	./bin/${APPNAME}

clean:
	rm -f ./bin/${APPNAME} ./cover.out

env:
	go env

sh:
	sh

zsh:
	zsh

stats:
	@echo ===============================================================================
	@echo Tests
	@echo ===============================================================================
	@cloc --match-f=_test.go ./src
	@echo
	@echo
	@echo ===============================================================================
	@echo Everything Else
	@echo ===============================================================================
	@cloc --not-match-f=_test.go ./src

# vendor
# * wipes our vendor directory completely
# * fetches all packages in ${PACKAGES} (GOPATH hardcoded to vendor directory only)
# * removes source control dirs so packages are vendorable (commitable)
vendor:
	rm -dRf ./${VENDOR}                              && \
	GOPATH=${PWD}/${VENDOR} go get -u ${PACKAGES}    && \
	rm -rf `find ./${VENDOR}/src -type d -name .git` && \
	rm -rf `find ./${VENDOR}/src -type d -name .hg`  && \
	rm -rf `find ./${VENDOR}/src -type d -name .bzr` && \
	rm -rf `find ./${VENDOR}/src -type d -name .svn`
