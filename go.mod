module github.com/cbarraford/cosmos-multisig

go 1.12

require (
	github.com/cbarraford/parsec v0.0.0-20190624083407-e967263e5f6b
	github.com/cosmos/cosmos-sdk v0.28.2-0.20190616100639-18415eedaf25
	github.com/gorilla/mux v1.7.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.0.3
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.31.5
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
