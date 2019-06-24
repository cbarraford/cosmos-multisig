# Building your app

## `Makefile`

Help users build your application by writing a `./Makefile` in the root directory that includes common commands:

> _*NOTE*_: The below Makefile contains some of same commands as the Cosmos SDK and Tendermint Makefiles.

```makefile
all: install

install: go.sum
    GO111MODULE=on go install -tags "$(build_tags)" ./cmd/nsd
    GO111MODULE=on go install -tags "$(build_tags)" ./cmd/nscli

go.sum: go.mod
    @echo "--> Ensure dependencies have not been modified"
    GO111MODULE=on go mod verify
```

### How about including Ledger Nano S support?

This requires a few small changes:

- Create a file `Makefile.ledger` with the following content:

```makefile
LEDGER_ENABLED ?= true

build_tags =
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
```

- Add `include Makefile.ledger` at the beginning of the Makefile:

```makefile
include Makefile.ledger

all: install

install: go.sum
    GO111MODULE=on go install -tags "$(build_tags)" ./cmd/nsd
    GO111MODULE=on go install -tags "$(build_tags)" ./cmd/nscli

go.sum: go.mod
    @echo "--> Ensure dependencies have not been modified"
    GO111MODULE=on go mod verify
```

## `go.mod`

Golang has a few dependency management tools. In this tutorial you will be using [`Go Modules`](https://github.com/golang/go/wiki/Modules). `Go Modules` uses a `go.mod` file in the root of the repository to define what dependencies the application needs. Cosmos SDK apps currently depend on specific versions of some libraries. The below manifest contains all the necessary versions. To get started replace the contents of the `./go.mod` file with the `constraints` and `overrides` below:

> _*NOTE*_: If you are following along in your own repo you will need to change the module path to reflect that (`github.com/{ .Username }/{ .Project.Repo }`).

```
module github.com/cosmos/sdk-application-tutorial

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.28.2-0.20190616100639-18415eedaf25
	github.com/gorilla/mux v1.7.0
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/prometheus/procfs v0.0.0-20190328153300-af7bedc223fb // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.0.3
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.31.5
	golang.org/x/sys v0.0.0-20190329044733-9eb1bfa1ce65 // indirect
	google.golang.org/genproto v0.0.0-20190327125643-d831d65fe17d // indirect
	google.golang.org/grpc v1.19.1 // indirect
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5

```

## Building the app

```bash
# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands:
nsd help
nscli help
```

### Congratulations, you have finished your nameservice application! Try [running and interacting with it](./build-run.md)!
