#!/bin/bash

# HTTPS (TLS) proxy address
PACHTTPSHOST='{{{PACHTTPSHOST}}}'

# Regular proxy address
PACPROXYHOST='{{{PACPROXYHOST}}}'

# Facebook and Twitter proxy address
PACFBTWHOST='{{{PACFBTWHOST}}}'

PACFILE="result/proxy-ssl.pac.mustache"
PACFILE_NOSSL="result/proxy.pac.mustache"

# Perform DNS resolving to detect and filter non-existent domains
RESOLVE_NXDOMAIN="no"
