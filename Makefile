BINARY_NAME=kubectl-meshsync_snapshot
VERSION=0.1.0
PLATFORMS=linux darwin windows
ARCHITECTURES=amd64

GO=go
GOFMT=gofmt
GOBUILD=$(GO) build
GOTEST=$(GO) test
GOGET=$(GO) get
GOMOD=$(GO) mod
RM=rm -f
MKDIR=mkdir -p
CHECKSUMS=sha256sum
GOLINT=golangci-lint run
TAR=tar
ZIP=zip
KREW=kubectl krew

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Set the correct binary suffix for Windows targets
ifeq ($(OS),Windows_NT)
    BINARY_SUFFIX=.exe
else
    BINARY_SUFFIX=
endif

all: build

.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME)$(BINARY_SUFFIX) cmd/kubectl-meshsync_snapshot/main.go

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: clean
clean:
	$(RM) $(BINARY_NAME)
	$(RM) $(BINARY_NAME).exe
	$(RM) -r dist/

.PHONY: lint
lint:
	$(GOLINT) ./...

.PHONY: fmt
fmt:
	$(GOFMT) -s -w .

.PHONY: tidy
tidy:
	$(GOMOD) tidy

.PHONY: deps
deps:
	$(GOGET) -u ./...

.PHONY: dist
dist: clean
	$(MKDIR) dist
	@for platform in $(PLATFORMS); do \
		for arch in $(ARCHITECTURES); do \
			echo "Building for $$platform/$$arch..."; \
			GOOS=$$platform GOARCH=$$arch $(GOBUILD) -o dist/$(BINARY_NAME)$(if $(filter windows,$$platform),.exe,) cmd/kubectl-meshsync_snapshot/main.go; \
			cp LICENSE dist/; \
			cd dist; \
			if [ "$$platform" = "windows" ]; then \
				if command -v $(ZIP) >/dev/null 2>&1; then \
					$(ZIP) -q $(BINARY_NAME)_$${platform}_$${arch}.zip $(BINARY_NAME).exe LICENSE; \
					$(CHECKSUMS) $(BINARY_NAME)_$${platform}_$${arch}.zip > $(BINARY_NAME)_$${platform}_$${arch}.zip.sha256; \
				else \
					echo "Warning: zip command not found, skipping Windows archive creation"; \
					$(CHECKSUMS) $(BINARY_NAME).exe > $(BINARY_NAME)_$${platform}_$${arch}.sha256; \
				fi; \
			else \
				$(TAR) -czf $(BINARY_NAME)_$${platform}_$${arch}.tar.gz $(BINARY_NAME) LICENSE; \
				$(CHECKSUMS) $(BINARY_NAME)_$${platform}_$${arch}.tar.gz > $(BINARY_NAME)_$${platform}_$${arch}.tar.gz.sha256; \
			fi; \
			cd ..; \
		done; \
	done

.PHONY: krew-install
krew-install: build
	cp $(BINARY_NAME)$(BINARY_SUFFIX) ~/.krew/bin/$(BINARY_NAME)$(BINARY_SUFFIX)

.PHONY: krew-validate
krew-validate:
	$(KREW) lint .krew.yaml

.PHONY: install
install: build
	cp $(BINARY_NAME)$(BINARY_SUFFIX) $(GOBIN)/$(BINARY_NAME)$(BINARY_SUFFIX)

.PHONY: uninstall
uninstall:
	$(RM) $(GOBIN)/$(BINARY_NAME)$(BINARY_SUFFIX)

.PHONY: help
help:
	@echo "Available make targets:"
	@echo "  all (default): build"
	@echo "  build:         Build the kubectl-meshsync-snapshot binary"
	@echo "  test:          Run unit tests"
	@echo "  clean:         Remove build artifacts"
	@echo "  lint:          Run linter checks"
	@echo "  fmt:           Format Go code"
	@echo "  tidy:          Tidy Go modules"
	@echo "  deps:          Update dependencies"
	@echo "  dist:          Build distribution binaries for all platforms"
	@echo "  krew-install:  Install plugin locally via Krew"
	@echo "  krew-validate: Validate the Krew manifest"
	@echo "  install:       Install plugin to GOBIN"
	@echo "  uninstall:     Remove plugin from GOBIN" 