LINTERS=\
	gofmt \
	golint \
	gosimple \
	vet \
	misspell \
	ineffassign \
	deadcode

ci: $(LINTERS) test

.PHONY: ci

#################################################
# Bootstrapping for base golang package deps
#################################################

BOOTSTRAP=\
	github.com/golang/lint/golint \
	honnef.co/go/tools/cmd/gosimple \
	github.com/client9/misspell/cmd/misspell \
	github.com/gordonklaus/ineffassign \
	github.com/tsenart/deadcode \
	github.com/alecthomas/gometalinter \
	github.com/Masterminds/glide

$(BOOTSTRAP):
	go get -u $@
bootstrap: $(BOOTSTRAP)
	cd ${GOPATH}/src/github.com/Masterminds/glide && git checkout -B v0.12.3 tags/v0.12.3 && go install && cd -

vendor: glide.lock
	glide install


.PHONY: bootstrap $(BOOTSTRAP)

#################################################
# Test and linting
#################################################

test: vendor
	@CGO_ENABLED=0 go test -v $(glide nv)

METALINT=gometalinter --tests --disable-all --vendor --deadline=5m -s data \
	 ./... --enable

$(LINTERS): vendor
	$(METALINT) $@

.PHONY: $(LINTERS) test
