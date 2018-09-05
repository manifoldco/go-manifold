ci: lint test

.PHONY: ci

#################################################
# Bootstrapping for base golang package deps
#################################################

BOOTSTRAP=\
  github.com/golangci/golangci-lint/cmd/golangci-lint \
	github.com/jbowes/oag

$(BOOTSTRAP):
	go get -u $@
bootstrap: $(BOOTSTRAP)

vendor: go.mod
	go get -v ./...

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

test: vendor $(GENERATED_NAMING_FILES)
	@CGO_ENABLED=0 go test -v $$(go list ./... | grep -v vendor)

lint:
	golangci-lint run --disable-all -E gofmt -E golint -E gosimple -E govet -E misspell -E ineffassign -E deadcode --skip-dirs=data

.PHONY: lint test

#################################################
# Releasing
#################################################

release:
ifneq ($(shell git rev-parse --abbrev-ref HEAD),master)
	$(error You are not on the master branch)
endif
ifneq ($(shell git status --porcelain),)
	$(error You have uncommitted changes on your branch)
endif
ifndef VERSION
	$(error You need to specify the version you want to tag)
endif
	sed -i -e 's|Version = ".*"|Version = "$(VERSION)"|' version.go
	git add version.go
	git commit -m "Tagging v$(VERSION)"
	git tag v$(VERSION)
	git push
	git push --tags


#################################################
# Data generation
#################################################

GO_BUILD=CGO_ENABLED=0 go build -i --ldflags="-w"

TOOLS=$(PREFIX)tools/bin

GENERATED_NAMING_FILES=$(patsubst names/data/%.txt,names/data/zz_generated_%.go,$(wildcard names/data/*.txt))
$(GENERATED_NAMING_FILES): names/data/zz_generated_%.go: $(TOOLS)/name-data names/data/%.txt
	$^ $@

TOOL_BINS=

define TOOL_BIN_TMPL
$(TOOLS)/$(1): vendor $$(call rwildcard,tools/$(1),*) $(2)
	$(3) $(GO_BUILD) -o $$@ ./tools/$(1)
TOOL_BINS += $(TOOLS)/$(1)
endef

$(eval $(call TOOL_BIN_TMPL,name-data))

tools: $(TOOL_BINS)

.PHONY: tools
