#!/usr/bin/env bash
RUN_NAME="magic.stock.fund"
mkdir -p output/conf
mkdir -p output/bin output/templates output/static
cp script/bootstrap.sh script/settings.py output
cp conf/* output/conf
cp -r static/* output/static 2>/dev/null
cp -r templates/* output/templates 2>/dev/null
chmod +x output/bootstrap.sh

go build -a -o output/bin/${RUN_NAME}