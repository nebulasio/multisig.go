VERSION?=1.0.7

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

CURRENT_DIR=$(shell pwd)
BUILD_DIR=${CURRENT_DIR}
BINARY=nebms

VET_REPORT=vet.report
LINT_REPORT=lint.report
TEST_REPORT=test.report
TEST_XUNIT_REPORT=test.report.xml

OS := $(shell uname -s)
ifeq ($(OS),Darwin)
	DYLIB=.dylib
	INSTALL=install
	LDCONFIG=
else
	DYLIB=.so
	INSTALL=sudo install
	LDCONFIG=sudo /sbin/ldconfig
endif

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH} -X main.compileAt=`date +%s`"

# Build the project
.PHONY: build build-linux clean dep lint run test vet link-libs

all: dep clean build

dep:
	dep ensure -v

build:
	make clean; go build $(LDFLAGS) -o $(BINARY)-$(COMMIT)
	ln -s $(BINARY)-$(COMMIT) $(BINARY)

clean:
	-rm -f $(BINARY)
	-rm -f $(BINARY)-$(COMMIT)
