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

TM_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')

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
			-X github.com/cometbft/cometbft/version.TMCoreSemVer=$(TM_VERSION)

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

all: install lint test

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

.PHONY: build build-linux install

###############################################################################
###                                Protobuf                                 ###
###############################################################################
protoVer=0.13.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-swagger:
	bash scripts/protoc-swagger-gen.sh

update-swagger-docs:
	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "Swagger docs are out of sync!";\
        exit 1;\
    else \
        echo "Swagger docs are in sync!";\
    fi

proto-all: proto-gen proto-swagger update-swagger-docs

.PHONY: proto-gen gen-swagger update-swagger-docs proto-all

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
	@go test -mod=readonly -timeout 10m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

benchmark:
	@go test -mod=readonly -bench=. ./...

simulate:
	@go test  -bench BenchmarkSimulation ./app -NumBlocks=200 -BlockSize 50 -Commit=true -Verbose=true -Enabled=true -Seed 1


test-e2e-pmf:
	cd interchaintest && go test -race -v -run TestPMF .

.PHONY: test test-all test-cover test-unit test-race simulate test-e2e-pmf

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	golangci-lint run --out-format=tab

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0

lint-docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run --timeout 10m

format-tools:
	go install mvdan.cc/gofumpt@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

format: format-tools
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*statik*" -not -name '*.pb.go' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*statik*" -not -name '*.pb.go' | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*statik*" -not -name '*.pb.go' | xargs goimports -w -local github.com/cosmos/cosmos-sdk

.PHONY: lint  lint-fix lint-docker format-tools format


###############################################################################
###                                Local Testnet (docker)                   ###
###############################################################################

localnet-rmi:
	$(DOCKER) rmi terra-money/localnet-core 2>/dev/null; true

localnet-build-env: localnet-rmi
	$(DOCKER) build --tag terra-money/localnet-core -f scripts/containers/Dockerfile \
			$(shell git rev-parse --show-toplevel)

localnet-build-nodes:
	$(DOCKER) run --rm -v $(CURDIR)/.testnets:/terra terra-money/localnet-core \
		testnet init-files --v 3 -o /terra --starting-ip-address 192.168.15.20 --keyring-backend=test --chain-id=core-testnet-1
	$(DOCKER) compose up -d

localnet-stop:
	$(DOCKER) compose down
	
localnet-start: localnet-stop localnet-build-env localnet-build-nodes

.PHONY: localnet-stop localnet-start localnet-build-env localnet-build-nodes localnet-rmi