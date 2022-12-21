#!/usr/bin/env bash
[[ -z "$1" ]] && echo "Usage: $0 input.nix" && exit 1
nix-instantiate --eval -E "(import ./$1).root"
