#!/bin/bash

set -e

project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
BUILD_DIR=$project_root/tests

TF_PROVIDER_VERSION="1.0.0-test"

setup_environment() {
    python3 --version || (echo "We need 'python3' installed to run integration tests"; exit 1)

    if [ ! -d "$project_root/.venv" ]; then
        echo "Creating virtual environment..."
        python3 -m venv $project_root/.venv
    fi

    source $project_root/.venv/bin/activate
    pip --version || (echo "We need 'pip' installed to run integration tests"; exit 1)

    echo "Installing dependencies..."
    pip install -q -r $BUILD_DIR/requirements.txt
}

build_provider() {
    if [ ! -f "$project_root/terraform-provider-outscale_v${TF_PROVIDER_VERSION}" ]; then
        echo "Building terraform provider version ${TF_PROVIDER_VERSION}..."
        cd $project_root
        make fmt
        make test
        go build -o terraform-provider-outscale_v${TF_PROVIDER_VERSION}
    fi
}

install_provider() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        PLUGIN_DIR=$BUILD_DIR/terraform.d/plugins/registry.terraform.io/outscale/outscale/${TF_PROVIDER_VERSION}/linux_amd64/
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        case $(uname -m) in
            arm64)
                PLUGIN_DIR=$BUILD_DIR/terraform.d/plugins/registry.terraform.io/outscale/outscale/${TF_PROVIDER_VERSION}/darwin_arm64/
                ;;
            *)
                PLUGIN_DIR=$BUILD_DIR/terraform.d/plugins/registry.terraform.io/outscale/outscale/${TF_PROVIDER_VERSION}/darwin_amd64/
                ;;
        esac
    else
        echo "OS $OSTYPE is not supported yet for testing"
        exit 1
    fi

    mkdir -p $PLUGIN_DIR
    cp $project_root/terraform-provider-outscale_v${TF_PROVIDER_VERSION} $PLUGIN_DIR/
}

run_tests() {
    local test_file=$1
    local parallel=$2

    cd $BUILD_DIR

    if [ -n "$parallel" ] && [ "$parallel" -gt 1 ] 2>/dev/null; then
        echo "Running $test_file with $parallel workers"
        python3 -m pytest -n $parallel -v $test_file || python3 -m pytest --lf -n $parallel -v $test_file
    else
        echo "Running $test_file sequentially"
        python3 -m pytest -v $test_file || python3 -m pytest --lf -v $test_file
    fi
}

cleanup() {
    rm -fr $BUILD_DIR/terraform.d
}
