#!/bin/bash

set -e

TOOL="${1:-terraform}"

if [ "$TOOL" != "terraform" ] && [ "$TOOL" != "tofu" ]; then
    echo "Usage: $0 [terraform|tofu]"
    echo "  Default: terraform"
    exit 1
fi

project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
EXAMPLES_DIR=$project_root/examples

echo "Running examples with $TOOL..."

for f in $EXAMPLES_DIR/*
do
    if [ -d $f ]
    then
        cd $f
        $TOOL init
        $TOOL validate
        cd -
    fi
done

exit 0
