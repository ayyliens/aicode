-include .env.default.properties
-include $(or $(CONF), .)/.env.properties

MAKEFLAGS := --silent --always-make
MAKE_PAR := $(MAKE) -j 128
WATCH := watchexec -r -c -d=0 -n
GOW := gow -c -v -w=go
VERB := $(if $(filter $(verb),true), -v,)
FAIL := $(if $(filter $(fail),false),,-failfast)
SHORT := $(if $(filter $(short),true), -short,)
CLEAR := $(if $(filter $(clear),false),,-c)
CLEAR_LONG := $(if $(filter $(clear),false),,--clear)
GO_SRC := ./go
GO_PKG := ./$(or $(pkg),$(GO_SRC)/...)
GO_FLAGS := -tags=$(tags) -mod=mod
GO_TEST_FLAGS := -count=1 $(GO_FLAGS) $(VERB) $(FAIL) $(SHORT)
GO_TEST_PATTERNS := -run="$(run)"

default:
	$(MAKE) go.run.w run=' \
		oai_conv_dir \
		--verb \
		--path=local/conv \
		--watch \
		--init \
		--funcs \
		$(CLEAR_LONG) \
	'

go.test.w:
	$(GOW) test $(GO_PKG) $(GO_TEST_FLAGS) $(GO_TEST_PATTERNS)

go.test:
	go test $(GO_PKG) $(GO_TEST_FLAGS) $(GO_TEST_PATTERNS)

go.run.w:
	$(GOW) run $(GO_FLAGS) $(GO_SRC) $(run)

go.run:
	go run $(GO_FLAGS) $(GO_SRC) $(run)

py.test.w:
	$(WATCH) -- $(MAKE) py.test

py.test:
	python3 test.py

# Assumes MacOS and Homebrew.
deps:
	brew install -q watchexec
