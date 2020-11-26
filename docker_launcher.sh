#!/bin/sh
COMMAND="privateness --data-dir $DATA_DIR --wallet-dir $WALLET_DIR $@"

adduser -D -u 10000 privateness

if [[ \! -d $DATA_DIR ]]; then
    mkdir -p $DATA_DIR
fi
if [[ \! -d $WALLET_DIR ]]; then
    mkdir -p $WALLET_DIR
fi

chown -R privateness:privateness $( realpath $DATA_DIR )
chown -R privateness:privateness $( realpath $WALLET_DIR )

su privateness -c "$COMMAND"
