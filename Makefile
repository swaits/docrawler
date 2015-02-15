# Golang project Makefile
#
# derived from: http://zduck.com/2014/go-project-structure-and-dependencies/

# build targets (which aren't files)
.PHONY: build fmt lint vet test autotest run clean env vendor

# configuration
APPNAME := docrawler
VENDOR := .vendor
PACKAGES :=                     \
	golang.org/x/tools/cmd/vet    \
	github.com/golang/lint/golint

# override GOPATH to be our vendor directory
export GOPATH := ${PWD}/${VENDOR}
export PATH := ${PWD}/${VENDOR}/bin:${PATH}

default: build

build: fmt lint vet
	go build -v -o ./bin/${APPNAME} ./src/${APPNAME}

fmt:
	go fmt ./src/...

lint:
	golint ./src

vet:
	go vet ./src/...

test:
	go test ./src/... | ./testfilter.sh

autotest:
	fswatch -or ./src | \
	xargs -n1 -I{} ${SHELL} -c 'echo ; date ; echo ----------------------------; go test ./src/... | ./testfilter.sh'

run: build
	./bin/${APPNAME}

clean:
	rm -f ./bin/${APPNAME}

env:
	go env

# vendor
# * wipes our vendor directory completely
# * fetches all packages in ${PACKAGES}
# * removes source control dirs so packages are vendorable (commitable)
vendor:
	rm -dRf ./${VENDOR}/src                          && \
	go get -u ${PACKAGES}                            && \
	rm -rf `find ./${VENDOR}/src -type d -name .git` && \
	rm -rf `find ./${VENDOR}/src -type d -name .hg`  && \
	rm -rf `find ./${VENDOR}/src -type d -name .bzr` && \
	rm -rf `find ./${VENDOR}/src -type d -name .svn`
