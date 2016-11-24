MAIN_PATH	:= github.com/sscheele/hallo

GLIDE		:= $(shell which glide)
GO		:= $(shell which go)
GO_BUILD_FLAGS	:=
GO_BUILD	:= $(GO) build $(GO_BUILD_FLAGS)
GO_TEST_FLAGS	:= -v
GO_TEST		:= $(GO) test $(GO_TEST_FLAGS)
GO_BENCH_FLAGS	:= -bench=. -benchmem
GO_BENCH	:= $(GO) test $(GO_BENCH_FLAGS)

TARGETS		:= bin/hallo

all: bin $(TARGETS)

bin:
	mkdir -p $@

bin/%: $(GLIDE) $(shell find . -name "*.go" -type f)
	$(GO_BUILD) -o $@ $(MAIN_PATH)

clean:
	-rm -f $(TARGETS)

test: $(GLIDE)
	$(GO_TEST) $(shell $(GLIDE) novendor)

bench: $(GLIDE)
	$(GO_BENCH) $(shell $(GLIDE) novendor)

.PHONY: clean test bench

