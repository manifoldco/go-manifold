LINTERS=$(shell grep "// lint" tools.go | awk '{gsub(/\"/, "", $$1); print $$1}' | awk -F / '{print $$NF}') \
	gofmt \
	vet \

ci: $(LINTERS) test

.PHONY: ci

#################################################
# Bootstrapping for base golang package and tool deps
#################################################

CMD_PKGS=$(shell grep '	"' tools.go | awk -F '"' '{print $$2}')

define VENDOR_BIN_TMPL
vendor/bin/$(notdir $(1)): vendor/$(1) | vendor
	go build -a -o $$@ ./vendor/$(1)
VENDOR_BINS += vendor/bin/$(notdir $(1))
vendor/$(1): go.sum
	GO111MODULE=on go mod vendor
endef

$(foreach cmd_pkg,$(CMD_PKGS),$(eval $(call VENDOR_BIN_TMPL,$(cmd_pkg))))

$(patsubst %,%-bin,$(filter-out gofmt vet,$(LINTERS))): %-bin: vendor/bin/%
gofmt-bin vet-bin:

vendor: go.sum
	GO111MODULE=on go mod vendor

mod-update:
	GO111MODULE=on go get -u -m
	GO111MODULE=on go mod tidy

mod-tidy:
	GO111MODULE=on go mod tidy

.PHONY: $(CMD_PKGS)
.PHONY: mod-update mod-tidy

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

$(LINTERS): %: vendor/bin/gometalinter %-bin vendor
	PATH=`pwd`/vendor/bin:$$PATH gometalinter --tests --disable-all --vendor \
		--deadline=5m -s data --skip generated --enable $@ ./...

.PHONY: $(LINTERS) test
.PHONY: cover all-cover.txt
cover: vendor $(GENERATED_NAMING_FILES)
	@CGO_ENABLED=0 go test -v -coverprofile=coverage.txt -covermode=atomic $$(go list ./... | grep -v vendor)


#################################################
# Releasing
#################################################

release: mod-tidy
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
