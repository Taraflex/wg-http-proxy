#!/bin/bash

# HTTPS (TLS) proxy address
PACHTTPSHOST='{{{PACHTTPSHOST}}}'

# Regular proxy address
PACPROXYHOST='{{{PACPROXYHOST}}}'

# Facebook and Twitter proxy address
PACFBTWHOST='{{{PACFBTWHOST}}}'

PACFILE="proxy-ssl.pac.mustache"
PACFILE_NOSSL="proxy.pac.mustache"

# Perform DNS resolving to detect and filter non-existent domains
RESOLVE_NXDOMAIN="no"
