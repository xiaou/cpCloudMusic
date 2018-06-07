NAME=cpCloudMusic
PKG=github.com/xiaou/cpCloudMusic/cmd
VERSION=0.0.2

LDFLAGS=-s -w -X main.Version=$(VERSION)

ifeq ($(OS),Windows_NT)
	OSFLAG=windows
	ARCH=386
	NAME := $(NAME)-$(OSFLAG)-$(ARCH).exe
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		OSFLAG=darwin
		ARCH=amd64
		NAME := $(NAME)-$(OSFLAG)-$(ARCH)
	endif
endif

BINDIR=bin

build: build_dir
	GOOS=$(OSFLAG) GOARCH=$(ARCH) CGO_ENABLED=1 \
		go build -ldflags "$(LDFLAGS)" \
		-o $(BINDIR)/$(NAME) $(PKG)

build_dir:
	@mkdir -p $(BINDIR)
	@rm -rf $(BINDIR)/*


.PHONY: build
