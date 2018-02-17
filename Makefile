dep = ${GOPATH}/bin/dep
curDir = $(shell pwd)
vendor = $(curDir)/vendor
sources = *.go

# Target
target = go-tradfri-server

# Development
devTarget = go-tradfri-server-dev
devDepepends = ${GOPATH}/src/github/moroen/*

all: $(target)

$(target): $(dep) $(vendor) $(sources) 
	go build -v

$(dep):
	go get -u github.com/golang/dep/cmd/dep

$(vendor):
	dep ensure -v

$(devTarget): $(sources) $(devDepends)
	rm -rf vendor
	go build -v -o $(devTarget)

dev: $(devTarget) 

test: dev
	./$(devTarget)

install: $(target)
	go install

clean:
	rm -rf $(vendor); rm -rf $(target)
