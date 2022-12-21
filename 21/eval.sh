#!/usr/bin/env bash
nix-instantiate --eval -E "(import ./$1).root"
