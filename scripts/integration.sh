#!/bin/bash

set -e
project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
BUILD_DIR=$project_root/tests/qa_provider_oapi
export TF_VAR_suffixe_lbu_name=$RANDOM

python3 --version || (echo "We need 'python3' intalled to run integration tests"; exit 1)
python3 -m venv .venv
source .venv/bin/activate
pip --version || (echo "We need 'pip' intalled to run integration tests"; exit 1)

make fmt
make test
go build -o terraform-provider-outscale_v0.5.32
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    mkdir -p $BUILD_DIR/terraform.d/plugins/registry.terraform.io/outscale/outscale/0.5.32/linux_amd64/
    cp terraform-provider-outscale_v0.5.32 $BUILD_DIR/terraform.d/plugins/registry.terraform.io/outscale/outscale/0.5.32/linux_amd64/
elif [[ "$OSTYPE" == "darwin"* ]]; then
    case $(uname -m) in
	arm64)
	    mkdir -p $BUILD_DIR/terraform.d/plugins/registry.terraform.io/outscale/outscale/0.5.32/darwin_arm64/
	    cp terraform-provider-outscale_v0.5.32 $BUILD_DIR/terraform.d/plugins/registry.terraform.io/outscale/outscale/0.5.32/darwin_arm64/
	    ;;
	*)
	    mkdir -p $BUILD_DIR/terraform.d/plugins/registry.terraform.io/outscale/outscale/0.5.32/darwin_amd64/
	    cp terraform-provider-outscale_v0.5.32 $BUILD_DIR/terraform.d/plugins/registry.terraform.io/outscale/outscale/0.5.32/darwin_amd64/
	    ;;
    esac
else
    echo "OS $OSTYPE is not supported yet for testing"
    exit 1
fi

cd $BUILD_DIR
pip install -r requirements.txt
pytest -v ./test_provider_oapi.py || pytest --lf -v ./test_provider_oapi.py
rm -fr terraform.d || exit 0
