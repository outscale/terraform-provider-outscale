#!/bin/bash

set -e 

DOC_TEMPLATE_SUBMODULE="docs/doc-terraform-template"

if [ "$TAG" == "" ]; then
    echo "We need the tag of the doc"
    exit 1
fi

# Create a branch
git checkout -b "autobuild-Documentation-$TAG"

# Update submodule
(cd $DOC_TEMPLATE_SUBMODULE && git fetch && git checkout $TAG && git submodule update)

# Gen the doc
make doc

# Create PR
git config user.name "Outscale Bot"
git config user.email "opensource+bot@outscale.com"

git add $DOC_TEMPLATE_SUBMODULE
git add docs/*

git commit -sm "Release Documentation $TAG"