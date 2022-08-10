#!/bin/sh

locked_version="$( cat .go-version )"
installed_version="$( go version )"

if echo "${installed_version}"  | grep -q "${locked_version}"; then
    exit 0
fi
echo "go version mismatch: installed go version does not match version in .go-version: ${locked_version}: installed version: ${installed_version}"
exit 1
