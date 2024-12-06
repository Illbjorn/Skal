#!/usr/bin/env bash

# Disable undeclared variable linting.
# shellcheck disable=SC2154

echo "Performing build for GOOS: ${GOOS}, GOARCH: ${GOARCH}."

# Create the build output directory if needed.
mkdir -p "${dir_bin}"

# Build the binary.
go build \
  -o "${dir_bin}/${name_bin}" \
  "${dir_app_main}"
