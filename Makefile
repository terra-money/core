#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
BUILDDIR ?= $(CURDIR)/build
DOCKER := $(shell which docker)
SHA256_CMD = sha256sum
GO_VERSION ?= "1.20"
# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --tags)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

TM_VERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::')

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += rocksdb
endif
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += boltdb
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=terra \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=terrad \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
			-X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TM_VERSION)

# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  $(info ################################################################)
  $(info To use rocksdb, you need to install rocksdb first)
  $(info Please follow this guide https://github.com/rockset/rocksdb-cloud/blob/master/INSTALL.md)
  $(info ################################################################)
  CGO_ENABLED=1
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=boltdb
endif

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# The below include contains the tools and runsim targets.
include contrib/devtools/Makefile

all: tools install lint test

build: go.sum
ifeq ($(OS),Windows_NT)
	exit 1
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/terrad ./cmd/terrad
endif

build/linux/amd64:
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o "$@/terrad" ./cmd/terrad

build/linux/arm64:
	GOOS=linux GOARCH=arm64 go build -mod=readonly $(BUILD_FLAGS) -o "$@/terrad" ./cmd/terrad

build/darwin/amd64:
	GOOS=darwin GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o "$@/terrad" ./cmd/terrad

build/darwin/arm64:
	GOOS=darwin GOARCH=arm64 go build -mod=readonly $(BUILD_FLAGS) -o "$@/terrad" ./cmd/terrad

build/windows/amd64:
	GOOS=windows GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o "$@/terrad" ./cmd/terrad

build-release: build/linux/amd64 build/linux/arm64 build/darwin/amd64 build/darwin/arm64 build/windows/amd64

build-linux:
	mkdir -p $(BUILDDIR)
	docker build --no-cache --tag terramoney/core ./
	docker create --name temp terramoney/core:latest
	docker cp temp:/usr/local/bin/terrad $(BUILDDIR)/
	docker rm temp

build-linux-with-shared-library:
	mkdir -p $(BUILDDIR)
	docker build --tag terramoney/core-shared ./ -f ./shared.Dockerfile
	docker create --name temp terramoney/core-shared:latest
	docker cp temp:/usr/local/bin/terrad $(BUILDDIR)/
	docker cp temp:/lib/libwasmvm.so $(BUILDDIR)/
	docker rm temp

build-release-amd64: go.sum $(BUILDDIR)/
	$(DOCKER) buildx create --name core-builder || true
	$(DOCKER) buildx use core-builder
	$(DOCKER) buildx build \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
    --build-arg BUILDPLATFORM=linux/amd64 \
    --build-arg GOOS=linux \
    --build-arg GOARCH=amd64 \
		-t core:local-amd64 \
		--load \
		-f Dockerfile .
	$(DOCKER) rm -f core-builder || true
	$(DOCKER) create -ti --name core-builder core:local-amd64
	$(DOCKER) cp core-builder:/usr/local/bin/terrad $(BUILDDIR)/release/terrad
	tar -czvf $(BUILDDIR)/release/terra_$(VERSION)_Linux_x86_64.tar.gz -C $(BUILDDIR)/release/ terrad
	rm $(BUILDDIR)/release/terrad
	$(DOCKER) rm -f core-builder

build-release-arm64: go.sum $(BUILDDIR)/
	$(DOCKER) buildx create --name core-builder  || true
	$(DOCKER) buildx use core-builder 
	$(DOCKER) buildx build \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
    --build-arg BUILDPLATFORM=linux/arm64 \
    --build-arg GOOS=linux \
    --build-arg GOARCH=arm64 \
		-t core:local-arm64 \
		--load \
		-f Dockerfile .
	$(DOCKER) rm -f core-builder || true
	$(DOCKER) create -ti --name core-builder core:local-arm64
	$(DOCKER) cp core-builder:/usr/local/bin/terrad $(BUILDDIR)/release/terrad 
	tar -czvf $(BUILDDIR)/release/terra_$(VERSION)_Linux_arm64.tar.gz -C $(BUILDDIR)/release/ terrad 
	rm $(BUILDDIR)/release/terrad
	$(DOCKER) rm -f core-builder

install: go.sum 
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/terrad

gen-swagger-docs:
	bash scripts/protoc-swagger-gen.sh

update-swagger-docs: statik
	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "Swagger docs are out of sync!";\
        exit 1;\
    else \
        echo "Swagger docs are in sync!";\
    fi

apply-swagger: gen-swagger-docs update-swagger-docs

.PHONY: build build-linux install update-swagger-docs apply-swagger


###############################################################################
###                        Integration Tests                                ###
###############################################################################

integration-test-all: init-test-framework \
	test-relayer \
	test-ica \
	test-ibc-hooks \
	test-vesting-accounts \
	test-alliance \
	test-tokenfactory
	-@rm -rf ./data
	-@killall terrad 2>/dev/null
	-@killall rly 2>/dev/null

init-test-framework: clean-testing-data install
	@echo "Initializing both blockchains..."
	./scripts/tests/start.sh

test-relayer:
	@echo "Testing relayer..."
	./scripts/tests/relayer/interchain-acc-config/rly-init.sh

test-ica: 
	@echo "Testing ica..."
	./scripts/tests/ica/delegate.sh

test-ibc-hooks: 
	@echo "Testing ibc hooks..."
	./scripts/tests/ibc-hooks/increment.sh

test-alliance: 
	@echo "Testing alliance module..."
	./scripts/tests/alliance/delegate.sh

test-vesting-accounts: 
	@echo "Testing vesting accounts..."
	./scripts/tests/vesting-accounts/validate-vesting.sh

test-tokenfactory: 
	@echo "Testing tokenfactory..."
	./scripts/tests/tokenfactory/tokenfactory.sh

clean-testing-data:
	@echo "Killing terrad and removing previous data"
	-@rm -rf ./data
	-@killall terrad 2>/dev/null
	-@killall rly 2>/dev/null

.PHONY: integration-test-all init-test-framework test-relayer test-ica test-ibc-hooks test-vesting-accounts test-tokenfactory clean-testing-data

###############################################################################
###                                Protobuf                                 ###
###############################################################################

proto-all: proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen:v0.3 sh ./scripts/protocgen.sh

.PHONY: proto-all proto-gen

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/terrad -d 2 | dot -Tpng -o dependency-graph.png

distclean: clean tools-clean
clean:
	rm -rf \
    $(BUILDDIR)/ \
    artifacts/ \
    tmp-swagger-gen/

.PHONY: distclean clean


###############################################################################
###                           Tests 
###############################################################################

test: test-unit

test-all: test-unit test-race test-cover

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./...

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

benchmark:
	@go test -mod=readonly -bench=. ./...

simulate:
	@go test  -bench BenchmarkSimulation ./app -NumBlocks=200 -BlockSize 50 -Commit=true -Verbose=true -Enabled=true -Seed 1

.PHONY: test test-all test-cover test-unit test-race simulate

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	golangci-lint run --out-format=tab

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0
.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -path "./tests/mocks/*" -not -name '*.pb.go' -not -path "./_build/*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -path "./tests/mocks/*" -not -name '*.pb.go' -not -path "./_build/*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -path "./tests/mocks/*" -not -name '*.pb.go' -not -path "./_build/*" | xargs goimports -w -local github.com/cosmos/cosmos-sdk
.PHONY: format
