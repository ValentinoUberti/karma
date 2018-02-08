#!/usr/bin/env bash
#check if HEAD is tagged
git describe --exact-match HEAD | grep -q fatal
grepRes=$?;
#if tagged then release to github
if [[ (( $grepRes )) ]]; 
then
    echo "Detected tag:"
    git describe --tags HEAD
    goreleaser --rm-dist
fi