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
	github.com/golang/dep/cmd/dep \
	github.com/alecthomas/gometalinter \
	github.com/jbowes/oag

$(BOOTSTRAP):
	go get -u $@
bootstrap: $(BOOTSTRAP)
	gometalinter --install

vendor: Gopkg.lock
	dep ensure


.PHONY: bootstrap $(BOOTSTRAP)

#################################################
# Code generation
#################################################

generated-%: specs/%.oag.yaml
	@oag -c $<
generated: $(patsubst specs/%.oag.yaml,generated-%,$(wildcard specs/*.oag.yaml))

.PHONY: generated

#################################################
# Test and linting
#################################################

test: vendor
	@CGO_ENABLED=0 go test -v $(go list ./... | grep -v vendor)

METALINT=gometalinter --tests --disable-all --vendor --deadline=5m -s data \
	 ./... --enable

$(LINTERS): vendor
	$(METALINT) $@

.PHONY: $(LINTERS) test

#################################################
# Releasing
#################################################

release:
ifneq ($(shell git rev-parse --abbrev-ref HEAD),master)
	$(error You are not on the master branch)
endif
ifndef VERSION
	$(error You need to specify the version you want to tag)
endif
	cat version.go | sed -e 's|Version = ".*"|Version = "$(VERSION)"|' > version.go
	git add .
	git commit -m "Tagging v$(VERSION)"
	git tag v$(VERSION)
	git push
	git push --tags
