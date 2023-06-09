#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
KEYRING_DIR="${SCRIPT_DIR}/keyring-test"
UNSIGNED_TX_PATH="${SCRIPT_DIR}/unsigned_tmp.json"
SIGNED_TX_PATH="${SCRIPT_DIR}/signed_tmp.json"
RUN1_PATH="${SCRIPT_DIR}/run1.json"
RUN2_PATH="${SCRIPT_DIR}/run2.json"
RUN3_PATH="${SCRIPT_DIR}/run3.json"
RUN4_PATH="${SCRIPT_DIR}/run4.json"
CHAIN_ID=example
ALICE=cosmos1swsy3qx89vqn2esx6fz4h9qym28u85hexkpkp9
BOB=cosmos1kq4jvak8d67x3204slu4jde85mmdpxpvyhaff9
CHRIS=cosmos1kk3g563ey8uljjm4falrw2q7f6pwcfznrn40tg
DARIA=cosmos1d8y6e8svzfcalwwqq8z22gdkcejul2mnk2swy6

rm $RUN1_PATH $RUN2_PATH $RUN3_PATH $RUN4_PATH 2> /dev/null

# Generate run #1
# it sends 1token from alice --> bob 5 times
hellod tx bank send $ALICE $BOB 1token --chain-id $CHAIN_ID --keyring-dir $KEYRING_DIR --generate-only > $UNSIGNED_TX_PATH 
for i in {0..5}
do
  hellod tx sign $UNSIGNED_TX_PATH --chain-id $CHAIN_ID --keyring-dir $KEYRING_DIR --from $ALICE --offline --sequence $i --account-number 1 > $SIGNED_TX_PATH
  hellod tx encode $SIGNED_TX_PATH  >> $RUN1_PATH
done

# Generate run #2
# it sends 1token from bob --> chris 100 times
hellod tx bank send $BOB $CHRIS 1token --chain-id $CHAIN_ID --keyring-dir $KEYRING_DIR --generate-only > $UNSIGNED_TX_PATH 
for i in {0..99}
do
  hellod tx sign $UNSIGNED_TX_PATH --chain-id $CHAIN_ID --keyring-dir $KEYRING_DIR --from $BOB --offline --sequence $i --account-number 2 > $SIGNED_TX_PATH
  hellod tx encode $SIGNED_TX_PATH  >> $RUN2_PATH
done

# Generate run #3
# it sends 1token from chris --> daria 10 times
hellod tx bank send $CHRIS $DARIA 1token --chain-id $CHAIN_ID --keyring-dir $KEYRING_DIR --generate-only > $UNSIGNED_TX_PATH 
for i in {0..9}
do
  hellod tx sign $UNSIGNED_TX_PATH --chain-id $CHAIN_ID --keyring-dir $KEYRING_DIR --from $CHRIS --offline --sequence $i --account-number 3 > $SIGNED_TX_PATH
  hellod tx encode $SIGNED_TX_PATH  >> $RUN3_PATH
done

# Generate run #4
# it sends 1token from daria --> alice 1 times
hellod tx bank send $DARIA $ALICE 1token --chain-id $CHAIN_ID --keyring-dir $KEYRING_DIR --generate-only > $UNSIGNED_TX_PATH 
for i in {0..1}
do
  hellod tx sign $UNSIGNED_TX_PATH --chain-id $CHAIN_ID --keyring-dir $KEYRING_DIR --from $DARIA --offline --sequence $i --account-number 4 > $SIGNED_TX_PATH
  hellod tx encode $SIGNED_TX_PATH  >> $RUN4_PATH
done