dep = ${GOPATH}/bin/dep
curDir = $(shell pwd)
# vendor = $(curDir)/vendor
target = go-tradfri-server
sources = *.go

all: $(target)

$(target): $(dep) $(vendor) $(sources) 
	go build -v

$(dep):
	go get -u github.com/golang/dep/cmd/dep

$(vendor):
	dep ensure -v

test: $(target)
	./$(target)

clean:
	rm -rf $(vendor); rm $(target)