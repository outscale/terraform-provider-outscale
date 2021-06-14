#!/bin/bash

set -e
project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
FILTER=$(ls $project_root | grep -v '^build$')

#Check tools
docker -v > /dev/null || ( echo "We need 'docker' intalled to run website locally"; exit 1)
make -v > /dev/null || (echo "We need 'make' intalled to run website locally"; exit 1)
git --version > /dev/null || (echo "We need 'git' intalled to run website locally"; exit 1)

if ! [ -d $project_root/build/doc/terraform-website ] ; then
    git clone https://github.com/hashicorp/terraform-website.git $project_root/build/doc/terraform-website
fi

BUILD_DIR=$project_root/build/doc/terraform-website
if ! [ -d $BUILD_DIR/ext/providers/outscale ] ; then
    mkdir $BUILD_DIR/ext/providers/outscale
fi

#Link and run the website
cp -r $FILTER $BUILD_DIR/ext/providers/outscale/
cd $BUILD_DIR/content/source/layouts/ &&
    ln -sf ../../../ext/providers/outscale/website/outscale.erb outscale.erb
cd $BUILD_DIR/content/source/docs/providers &&
    ln -sf ../../../../ext/providers/outscale/website/docs outscale
cd $BUILD_DIR  && make website
