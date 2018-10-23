.PHONY: build
build: $(GOPATH)/bin/stringer
	go generate
	go build

precommit:
	gofmt -w .
	go test

$(GOPATH)/bin/stringer:
	go get golang.org/x/tools/cmd/stringer
