#!/bin/bash

set -e

project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
EXAMPLES_DIR=$project_root/examples
INSTALL_DIR=terraform.d/plugins/registry.terraform.io/outscale-dev/outscale/0.5.32/linux_amd64/

go build -o terraform-provider-outscale_v0.5.32

for f in $EXAMPLES_DIR/*
do
    if [ -d $f ]
    then
        cd $f
        mkdir -p $f/$INSTALL_DIR
        cp ../../terraform-provider-outscale_v0.5.32 $f/$INSTALL_DIR
        terraform init
        terraform apply -auto-approve
        terraform destroy -auto-approve
        cd -
    fi
done

exit 0