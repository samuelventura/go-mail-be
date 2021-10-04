#!/bin/bash -x

if [[ "$OSTYPE" == "linux"* ]]; then
    SRC=$HOME/go/bin
    DST=/usr/local/bin
    if [[ -f "$DST/go-mail-ss" ]]; then
        sudo systemctl stop GoMailMs
        sudo $DST/go-mail-ss -service uninstall
        sleep 3
    fi
    go install
    (cd go-mail-ss; go install)
    sudo cp $SRC/go-mail-ms $DST
    sudo cp $SRC/go-mail-ss $DST
    sudo $DST/go-mail-ss -service install
    sudo systemctl restart GoMailMs
fi
