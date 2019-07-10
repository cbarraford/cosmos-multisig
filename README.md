Cosmos Multi-Signature Wallet
=============================
A custom Cosmos SDK blockchain that supports multi-signature wallets.

## Setup
Ensure you have a recent version of go (ie `1.121) and enabled go modules
```
export GO111MODULE=on
```
And have `GOBIN` in your `PATH`
```
export GOBIN=$GOPATH/bin
```

### Install
Install via this `make` command.

```bash
make install
```

Once you've installed `msigcli` and `msigd`, check that they are there.

```bash
msigcli help
msigd help
```

### Configuration

Next configure your chain.
```bash
# Initialize configuration files and genesis file
# moniker is the name of your node
msigd init <moniker> --chain-id msigchain


# Copy the Address output here and save it for later use
# [optional] add "--ledger" at the end to use a Ledger Nano S
msigcli keys add jack

# Copy the Address output here and save it for later use
msigcli keys add alice

# Add both accounts, with coins to the genesis file
msigd add-genesis-account $(msigcli keys show jack -a) 1000msig,100000000stake
msigd add-genesis-account $(msigcli keys show alice -a) 1000msig,100000000stake

# Configure your CLI to eliminate need for chain-id flag
msigcli config chain-id msigchain
msigcli config output json
msigcli config indent true
msigcli config trust-node true

msigd gentx --name jack
```

### Start
There are three services you may want to start.

#### Daemon
This runs the backend
```bash
msigd start
```

#### API Service
Starts an HTTP service to service requests to the backend.
```bash
msigcli rest-server
```

#### CORS Proxy
For making requests in a browser to the API backend, you'll need to start a
proxy in front of the API service to give proper CORS headers. 
For development purposes, a service is provided in `/scripts/cors`

```bash
cd scripts/cors
npm start
```
