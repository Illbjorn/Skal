#!/usr/bin/env bash

files=$(find . -type f -name '*.go')

total=0
for f in $files; do
  lines=$(wc -l "${f}" | awk '{print $1}')
  ((total+=lines))
done

echo $total
