#!/bin/bash

set -x
set -e

while true; do

  appHome="$HOME/go/src/github.com"
  appPath="cbarraford/cosmos-multisig"

  cd ${appHome}/${appPath}
  make install
  msigd init local --chain-id msigchain

  msigcli keys add jack
  msigcli keys add alice

  msigd add-genesis-account $(msigcli keys show jack -a) 1000atom,100000000stake
  msigd add-genesis-account $(msigcli keys show alice -a) 1000atom,100000000stake

  msigcli config chain-id msigchain
  msigcli config output json
  msigcli config indent true
  msigcli config trust-node true

  echo "password" | msigd gentx --name jack
  msigd collect-gentxs
  msigd validate-genesis

  msigd start & msigcli rest-server --chain-id msigchain --trust-node && fg

  break

done
