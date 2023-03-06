GO ?= go
EXECUTABLE := codegpt
GOFILES := $(shell find . -type f -name "*.go")
TAGS ?=
LDFLAGS ?= -X 'github.com/appleboy/CodeGPT/cmd.Version=$(VERSION)' -X 'github.com/appleboy/CodeGPT/cmd.Commit=$(COMMIT)'

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

ifneq ($(DRONE_TAG),)
	VERSION ?= $(DRONE_TAG)
else
	VERSION ?= $(shell git describe --tags --always || git rev-parse --short HEAD)
endif
COMMIT ?= $(shell git rev-parse --short HEAD)

build: $(EXECUTABLE)

$(EXECUTABLE): $(GOFILES)
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o $@

install: $(GOFILES)
	$(GO) install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)'

test:
	@$(GO) test -v -cover -coverprofile coverage.txt ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1
