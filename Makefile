### repo name
REPO_PATH := "github.com/abyss414/house"

### go path
GOPATHVAR=${GOPATH}

### bash
SHELL=/usr/bin/env bash
%.c: %.y
%.c: %.l

### version,branch,commit,date,changes,
VERSION := ""

PACK_SRV := ""


VERSIONCMD = "`git describe --exact-match --tags $(git log -n1 --pretty='%h')`"
VERSION := $(shell echo $(VERSIONCMD))
ifeq ($(strip $(VERSION)),)
  BRANCHCMD := "`git describe --contains --all HEAD`-`git rev-parse HEAD`"
  VERSION = $(shell echo $(BRANCHCMD))
else
  TAGCMD := "`git describe --exact-match --tags $(git log -n1 --pretty='%h')`-`git rev-parse HEAD`"
  VERSION =  $(shell echo $(TAGCMD))
endif
VERSION ?= $(VERSION)

BRANCHCMD := "`git describe --contains --all HEAD`"
BRANCH := $(shell echo $(BRANCHCMD))
BRANCH  ?= $(BRANCH)
COMMITCMD = "`git rev-parse HEAD`"
COMMIT := $(shell echo $(COMMITCMD))
DATE := $(shell echo `date +%FT%T%z`)
CHANGES := $(shell echo `git status --porcelain | wc -l`)
ifneq ($(strip $(CHANGES)), 0)
       VERSION := dirty-build-$(VERSION)
       COMMIT := dirty-build-$(COMMIT)
endif

REMOVESYMBOL := -w -s
ifeq (true, $(DEBUG))
       REMOVESYMBOL =
       GCFLAGS=-gcflags=all="-N -l "
endif

### CLDFLAGS,LDFLAG define
ifdef IS_LINUX
       CLDFLAGS += -L/usr/local/yay/lib
else ifdef IS_MAC_OS_X
       CLDFLAGS += -L /usr/local/lib/gcc/4.9 -L/usr/local/lib
endif


### build dir
BUILD_DIR := $(CURDIR)/build
PACKAGE_DIR := $(CURDIR)/package
BUILD := $(BUILD_DIR)/built
PROJECT_DIR = $(CURDIR)
BUILD_TAGS += gm no_development

export GOBIN := $(BUILD_DIR)/bin

### build step
default: pack

pack:
	@make build
	@make push

build:
	@make build-house_scrapper
	@make build-house_worker


build-house_scrapper:
	command docker build . --build-arg PACK_SRV="house_scrapper" --tag ulysseskk/house_scrapper:v$(VERSION)

build-house_worker:
	command docker build . --build-arg PACK_SRV="house_worker" --tag ulysseskk/house_worker:v$(VERSION)


push:
	docker push ulysseskk/house_worker:v$(VERSION)
	docker push ulysseskk/house_scrapper:v$(VERSION)
