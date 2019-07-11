Cosmos Multi-Signature Wallet
=============================
A custom Cosmos SDK blockchain that supports multi-signature wallets.

# Overview
While multsig transfer of coins utilizes the built in multisig feautres of
cosmos, this blockchain expands that to create an infrastructure to have a
nice UI around multisig transactions.

The complete workflow of a multisig transaction as follows...
 1. Create a multisig wallet, specifying the name, minimum number of
    signatures for a transaction to take place, and the pub keys associated
with the wallet.
 2. Transfer funds to this wallet
 3. Create a transaction request to send funds from the multisig wallet to
    another address
 4. Multiple owners of the multisig wallet sign the transaction.
 5. Once enough signatures have been made, send the funds.
 6. Update the transaction request with the `txhash` of the transfer of funds.

### CLI
The cli tool has a series of queries and transaction that you can use. In
theory, you could iteract with this blockchain fully using the cli, but REST
was the intended purpose.

#### Create a wallet
This command is used to create a multisig wallet. Pubkeys should be comma
separated (no spaces).
```
msgicli tx multisig create-wallet [name] [min-signatures-required] [pub-keys], [addresses] [flags]
```

#### Get a wallet
Get wallet info by wallet address
```
msgicli query multisig get-wallet [address] [flags]
```

#### Query wallets
Search for a list of wallets by public key
```
msgicli query multisig query-wallets [pub_key] [flags]
```

#### Create a transaction
This command creates a transaction request to move funds out of a multisig
wallet.
```
msgicli tx multisig create-transaction [from] [to] [coins] [signers] [flags]
```

#### Get transaction
Retrieve transaction request information by uuid
```
msgicli query multisig get-transaction [uuid] [flags]
```

#### Query transactions
Get a list of transaction requests by wallet address.
```
msgicli query multisig query-transactions [wallet_address] [flags]
```

#### Add signature to transaction
This command adds a signature to a transaction request.
TODO: remove need to supply `pubkey_base64`. This info is available via the
account info (`/auth/accounts/<address>`). 
```
msgicli tx multisig save-transaction-signature [uuid] [pubkey] [pubkey_base64] [signature] [signers] [flags]
```

#### Add TxHash to transaction
Once the transaction is completed and funds sent, save the `txhash` in the
transaction request to mark it as completed.
```
msgicli tx multisig complete-transaction [uuid] [transaction_id] [signers] [flags]
```

### API
There are corresponding API endpoints for each of the CLI commands above.

#### `POST /multisig/wallet`
Create a wallet

```
{
    "name": "demo 1",
    "base_req": {"chain_id":"msigchain", "from": "msigXXXXXX"},
    "min_sig_tx": 2,
    "pub_keys": [...],
    "signers": [...]
}
```

#### `GET /multisig/wallet/<address>`
Get a wallet

#### `GET /multisig/wallets/<pubkey>`
List wallets that contain specified public key

#### `POST /multisig/transaction`
Create a transaction request

```
{
    "base_req": {"chain_id":"msigchain", "from": "msigXXXXXX"},
    "from": "msigXXXX",
    "to": "msigXXXX",
    "amount": 3,
    "denom": "msigtoken",
    "signers": [...]
}
```

#### `GET /multisig/transaction/<uuid>`
Get a transaction request by uuid

#### `GET /multisig/transactions/<address>`
List transaction by wallet address

#### `POST /multisig/transaction/sign`
Add signature for a transaction request

```
{
    "base_req": {"chain_id":"msigchain", "from": "msigXXXXXX"},
    "uuid": "02206ab8-ef05-4ecc-8e81-4430405e929a",
    "signature": "1KRP93NJ85SxygGucoS7MqV39INDG/TYZzRP9NVzS4k/J7zn5kus1r3SoIXkyHgYEZmnaIyr26liL7uez45Uiw",
    "pub_key": "msigp1addwnpepq0wx75hs3jpepjlvvee4r8gmuuxnmjzk6k5jkjps7n9rr3y4v0quqqy456c",
    "pub_key_base64": "A/HRtDdV5mtPOMFhDDRRqb0s60q8rVJ+3AqSWd3PkOWx",
    "signers": [...]
}
```

#### `POST /multisig/sign/multi`
With given signatures, generate a multi-signature. 
"Signatures" must be a list of tx signatures for pub keys of the wallet. Order
is important here, and must align with the pub key order of the wallet.
"Slots" is a string of zeros and ones representing which pubkeys of the wallet
are included in the list of signatures, and which are not. Zeros are not
include, ones are included.
The resulting json response will include the multisig signature.

```
{
    "Signatures": [...],
    "Slots": "011",
}
```

#### `POST /multisig/transaction/complete`
Complete a transaction supplying the `txhash` of the transfer of funds.

```
{
    "base_req": {"chain_id":"msigchain", "from": "msigXXXXXX"},
    "uuid": "02206ab8-ef05-4ecc-8e81-4430405e929a",
    "TxID": "939HDJ300...",
    "signers": [...],
}
```

#### `POST /multisig/broadcast`
Broadcast a message (same as to `/txs` in the cosmos SDK).

# Developer

## Types

### `MultiSigWallet`
`MultiSigWallet` is a type to store multisignature wallet information. This
type includes...
 * `Name` - the custom name of the wallet. This makes it easier for users to
   identify each wallet and its purpose.
 * `MinSigTx` - the minimum number of regular user signatures required before
   a transaction can be sent.
 * `Address` - The receiving address to send coins into this wallet.
 * `PubKeys` - A list of public keys associated with this wallet that has the
   ability to sign transactions. Order of public keys is important.

** Notes ** Wallets cannot be deleted, nor can they be overwritten or change
once created.

### `Transaction`
`Transaction` is a type to store a transaction request information to move
funds out of a multisig wallet. 
 * `UUID` - a unique identifier (follow uuid standards) 
 * `From` - an multisig wallet address to send the funds from
 * `To` - a wallet address to send the funds to
 * `Coins` - an array of coins to be sent from the multisig wallet. Currently
   only one coins (one denom) can be sent at this time.
 * `Signatures` - the signed signatures of this transaction from the public
   keys associated with the multisig wallet
 * `TxID` - the transaction hash from the blockchain referencing this
   transaction on the blockchain. This is written as a last step to signify
the transaction is complete.
 * `CreatedAt` - The block height when this transaction request was first
   created. This helps the UI sort the transaction list, but also acts a means
to cleanup old transaction requests from history (ie deleting transaction
requests after X blocks have passed).

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
