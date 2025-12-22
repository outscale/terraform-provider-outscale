#!/bin/bash

set -e

project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
EXAMPLES_DIR=$project_root/examples

go build -o terraform-provider-outscale

for f in $EXAMPLES_DIR/*
do
    if [ -d $f ]
    then
        cd $f
        VERSION_DIR=`grep -o '[[:digit:]]*\.[[:digit:]]*\.[[:digit:]]*' outscale.tf`
        INSTALL_DIR=$f/terraform.d/plugins/registry.opentofu.org/outscale/outscale/$VERSION_DIR/linux_amd64/
        echo $INSTALL_DIR
        mkdir -p $INSTALL_DIR
        cp ../../terraform-provider-outscale $INSTALL_DIR
        tofu init
        tofu apply -auto-approve
        tofu destroy -auto-approve
        cd -
    fi
done

exit 0
