#!/bin/bash

set -e

cp -f -T .github/workflows/config.sh ./antizapret-pac-generator-light/config/config.sh

export RESOLVE_NXDOMAIN="${RESOLVE_NXDOMAIN:-yes}"

mkdir -p ./result
cd ./antizapret-pac-generator-light
sudo update-alternatives --set awk $(update-alternatives --list awk | grep gawk)
#echo > ./config/exclude-hosts-dist.txt
sed -i 's/\\\_/_/' parse.sh
pip install dnspython
bash ./doall.sh
bash ./generate-pac.sh

source config/config.sh

PACFILE_NOSSL_MIN="${PACFILE_NOSSL%.*}.min.js"

npx terser $PACFILE_NOSSL --ecma 5 --toplevel --mangle reserved=[FindProxyForURL,url,host] -o $PACFILE_NOSSL_MIN
sed -i "1s/^/\/\/ Generated on $(date --utc)\n/" $PACFILE_NOSSL_MIN

sed -i 's+function FindProxyForURL(url, host) {+function FindProxyForURL(url, host) {\n  {{#PACPROXYRE}} if (/{{{PACPROXYRE}}}/i.test(host)) return "PROXY {{{PACPROXYHOST}}}; DIRECT"; {{/PACPROXYRE}}\n  {{#PACDIRECTRE}} if (/{{{PACDIRECTRE}}}/i.test(host)) return "DIRECT"; {{/PACDIRECTRE}}\n+' $PACFILE_NOSSL
sed -i 's+function FindProxyForURL(url,host){+function FindProxyForURL(url, host) {\n  {{#PACPROXYRE}} if (/{{{PACPROXYRE}}}/i.test(host)) return "PROXY {{{PACPROXYHOST}}}; DIRECT"; {{/PACPROXYRE}}\n  {{#PACDIRECTRE}} if (/{{{PACDIRECTRE}}}/i.test(host)) return "DIRECT"; {{/PACDIRECTRE}}\n+' $PACFILE_NOSSL_MIN

brotli -kf $PACFILE_NOSSL_MIN -o $PACFILE_NOSSL_MIN.br
brotli -kf $PACFILE_NOSSL -o $PACFILE_NOSSL.br