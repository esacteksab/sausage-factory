#!/usr/bin/env bash

# User Data Scripts
export UDS
UDS="$(find . -name '*tpl' -type f -depth 2|cut -d'/' -f2|sort -u)"

# Branch Name Prefix
export BNP
BNP="bm-nt-cwlg-uds-"

for file in ${UDS}; do  git branch ${BNP}${file} main; done
