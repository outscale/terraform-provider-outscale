#!/bin/bash

set -e
project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
EXAMPLES_DIR=$project_root/examples

go build -o terraform-provider-outscale_v0.5.32
mkdir -p $BUILD_DIR/terraform.d/plugins/registry.terraform.io/outscale-dev/outscale/0.5.32/linux_amd64/
cp terraform-provider-outscale_v0.5.32 $BUILD_DIR/terraform.d/plugins/registry.terraform.io/outscale-dev/outscale/0.5.32/linux_amd64/

for f in $EXAMPLES_DIR/*
do
    if [ -d $f ]
    then
        cd $f
        terraform init
        terraform apply -auto-approve
        terraform destroy -auto-approve
        cd -
    fi
done

exit 0