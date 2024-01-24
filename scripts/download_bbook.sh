#!/bin/bash

mkdir -p bin

if [ ! -f bin/bbook ]
then
    echo "downloading banana-book"
    curl -sSL https://github.com/dfirebaugh/bbook/releases/download/v0.0.0/bbook-x86_64_unknown-linux.tar.gz | tar -xz --directory=bin
else
    echo "bin/bbook already exists"
fi
