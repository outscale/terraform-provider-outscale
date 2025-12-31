#!/bin/bash

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <test_file>"
    echo "Environment variables:"
    echo "  PYTEST_PARALLEL - Number of parallel workers (running sequentially if not set)"
    echo "  RUN_NETS_ONLY - Only run nets tests"
    echo "  SKIP_NETS - Skip nets tests"
    exit 1
fi

TEST_FILE=$1

script_dir=$(cd "$(dirname $0)" && pwd)
source "$script_dir/test-integration-common.sh"

setup_environment
build_provider
install_provider

run_tests $TEST_FILE $PYTEST_PARALLEL

cleanup
