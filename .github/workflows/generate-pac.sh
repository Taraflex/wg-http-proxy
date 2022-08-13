#!/bin/bash

set -e

cp -f -T .github/workflows/config.sh ./antizapret-pac-generator-light/config/config.sh
mkdir -p ./result
cd ./antizapret-pac-generator-light
sudo update-alternatives --set awk $(update-alternatives --list awk | grep gawk)
echo > ./config/exclude-hosts-dist.txt
sed -i 's/\\\_/_/' parse.sh
bash ./doall.sh
bash ./generate-pac.sh

source config/config.sh
brotli -kf $PACFILE_NOSSL -o $PACFILE_NOSSL.br