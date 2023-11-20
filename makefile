# For local use in this makefile. This does not export to sub-processes.
-include .env.default.properties
-include $(or $(CONF), .)/.env.properties

MAKEFLAGS := --silent --always-make
MAKE_PAR := $(MAKE) -j 128
VERB_SHORT := $(if $(filter $(verb),true), -v,)
VERB_LONG := $(if $(filter $(verb),false),,--verb)
CLEAR_SHORT := $(if $(filter $(clear),false),,-c)
CLEAR_LONG := $(if $(filter $(clear),true),--clear,)
GO_SRC := ./go
GO_PKG := ./$(or $(pkg),$(GO_SRC)/...)
GO_FLAGS := -tags=$(tags) -mod=mod
GO_RUN_ARGS := $(GO_FLAGS) $(GO_SRC) $(run)
GO_TEST_FAIL := $(if $(filter $(fail),false),,-failfast)
GO_TEST_SHORT := $(if $(filter $(short),true), -short,)
GO_TEST_FLAGS := -count=1 $(GO_FLAGS) $(VERB_SHORT) $(GO_TEST_FAIL) $(GO_TEST_SHORT)
GO_TEST_PATTERNS := -run="$(run)"
GO_TEST_ARGS := $(GO_PKG) $(GO_TEST_FLAGS) $(GO_TEST_PATTERNS)
OPEN_AI_API_AUTH := "Authorization: Bearer $(OPEN_AI_API_KEY)"

# Dependency: https://github.com/mitranim/gow.
GOW := gow $(CLEAR_SHORT) $(VERB_SHORT)

# Dependency: https://github.com/mattgreen/watchexec.
WATCH := watchexec $(CLEAR_SHORT) -d=0 -r -n

# TODO: if appropriate executable does not exist, print install instructions.
ifeq ($(OS),Windows_NT)
	GO_WATCH := $(WATCH) -w=$(GO_SRC) -- go
else
	GO_WATCH := $(GOW) -w=$(GO_SRC)
endif

ifeq ($(OS),Windows_NT)
	RM_DIR = if exist "$(1)" rmdir /s /q "$(1)"
else
	RM_DIR = rm -rf "$(1)"
endif

ifeq ($(OS),Windows_NT)
	CP_INNER = if exist "$(1)" copy "$(1)"\* "$(2)" >nul
else
	CP_INNER = if [ -d "$(1)" ]; then cp -r "$(1)"/* "$(2)" ; fi
endif

ifeq ($(OS),Windows_NT)
	CP_DIR = if exist "$(1)" copy "$(1)" "$(2)" >nul
else
	CP_DIR = if [ -d "$(1)" ]; then cp -r "$(1)" "$(2)" ; fi
endif

default:
	$(MAKE) go.run.w run=' \
		oai_conv_dir \
		$(VERB_LONG) \
		--path=local/conv \
		--watch \
		--init \
		--funcs \
		$(CLEAR_LONG) \
	'

example.xln:
	$(eval TAR := "local/conv_xln")
	$(call RM_DIR,$(TAR))
	$(call CP_DIR,local_example/conv_example_xln,$(TAR))

	$(MAKE) go.run.w run=' \
		oai_conv_dir \
		$(VERB_LONG) \
		--path=$(TAR) \
		--funcs \
		--src-path=$(TAR)/src_files \
		--out-path=$(TAR)/out_files \
		--watch \
		--trunc \
		--fork \
		--init \
		$(CLEAR_LONG) \
	'

go.run.w:
	$(GO_WATCH) run $(GO_RUN_ARGS)

go.run:
	go run $(GO_RUN_ARGS)

go.test.w:
	$(GO_WATCH) test $(GO_TEST_ARGS)

go.test:
	go test $(GO_TEST_ARGS)


open_ai_api_models:
	curl https://api.openai.com/v1/models -H $(OPEN_AI_API_AUTH)
	