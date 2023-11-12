VERSION ?= $(shell git describe --abbrev=7 || echo -n "unversioned")
VERSION_PACKAGE ?= github.com/pdcalado/gomutate/version

LDFLAGS ?= "-X '$(VERSION_PACKAGE).Version=$(VERSION)' -s -w"

GOBUILD ?= GCO_ENABLED=0 go build -ldflags=$(LDFLAGS) -tags osusergo,netgo

fmt:
	gofmt -w -s ./
	goimports -w -local github.com/pdcalado/gomutate ./

lint:
	golangci-lint run -v

clean:
	rm -rf ./bin

build:
	$(GOBUILD) -o gomutate .

test:
	go run main.go -type Acme ./testdata/acme.go > testdata/mutations.go
	go run testdata/*.go | diff - testdata/expected.txt
