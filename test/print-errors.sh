#!/bin/bash -e

suite=$1

pattern="$(dirname "$0")/suites/$suite/*.log"
for file in $pattern; do
    echo -e "\nFiltered errors from $file"
    cat $file | grep --ignore-case --extended-regexp --word-regexp "error|ERR"
done
