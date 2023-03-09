#!/bin/bash -e

suite=$1

pattern="$(dirname "$0")/suites/$suite/*.log"
for file in $pattern; do
    echo -e "\n🟥 Filtered errors from $file:"
    cat $file | grep --ignore-case --extended-regexp --word-regexp "error|ERR"
done

echo -e "\n🔎 For full logs, refer to workflow artifacts."
