#!/bin/sh

set -e

CURRENT_VERSION=$(cat ".go-version")
if [ -z $CURRENT_VERSION ] ; then
	echo "no current version found"
	exit 1
fi

NEW_VERSION=$1
if [ -z $1 ] ; then
	echo "no new version specified as parameter"
	exit 1
fi

echo "Updating go version from $CURRENT_VERSION to $NEW_VERSION"

# update all references (CI, docker)
for f in $(find . -type f -not -path '*testdata*'); do
	sed -i "s/go$CURRENT_VERSION/go$NEW_VERSION/g" $f
	sed -i "s/go:$CURRENT_VERSION/go:$NEW_VERSION/g" $f
	sed -i "s/golang:$CURRENT_VERSION/golang:$NEW_VERSION/g" $f
done

# update go version used in documentation
for f in $(find . build -type f -name 'md'); do
    sed -i "s/$CURRENT_VERSION/$NEW_VERSION/g" README.md
done

# finally, update .go-version
echo $NEW_VERSION > .go-version
