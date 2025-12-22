#!/bin/bash

set -e
project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
BUILD_DIR=$project_root/tests
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

if [ -n "$RUN_NETS_ONLY" ]; then
    PARALLEL_VALUE=$PYTEST_NETS_PARALLEL
else
    PARALLEL_VALUE=$PYTEST_PARALLEL
fi

if [ -n "$PARALLEL_VALUE" ] && [ "$PARALLEL_VALUE" -gt 1 ] 2>/dev/null; then
    echo "Running tests with $PARALLEL_VALUE workers"
    pytest -n $PARALLEL_VALUE -v ./test_provider_oapi.py || pytest --lf -n $PARALLEL_VALUE -v ./test_provider_oapi.py
else
    echo "Running tests sequentially"
    pytest -v ./test_provider_oapi.py || pytest --lf -v ./test_provider_oapi.py
fi

rm -fr terraform.d || exit 0
