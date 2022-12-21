#!/usr/bin/env bash
[[ -z "$1" ]] && echo "Usage: $0 input.qalc" && exit 1
qalc -f "$1" -nocurrencies -nodatasets -nounits -novariables -u8
