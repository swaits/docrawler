# Golang project Makefile
#
# derived from: http://zduck.com/2014/go-project-structure-and-dependencies/

# build targets (which aren't files)
.PHONY: build fmt lint vet test autotest cover doc run clean env vendor

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

lint:
	# TODO golint is not accepting the "./src/..." for some reason (from within make)
	golint ./src/${APPNAME}/*.go

vet:
	go vet ./src/...

test:
	go test -v ./src/... | ./util/testfilter.sh

autotest:
	fswatch -or ./src | \
		xargs -n1 -I{} ${SHELL} -c 'echo ; date ; echo ----------------------------; go test ./src/... | ./util/testfilter.sh'

cover:
	go test -coverprofile cover.out ./src/...  && go tool cover -html cover.out

doc:
	godoc -http=:6060 -index

run: build
	./bin/${APPNAME}

clean:
	rm -f ./bin/${APPNAME} ./cover.out

env:
	go env

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
