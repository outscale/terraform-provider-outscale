#!/bin/bash

set -e
project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
BUILD_DIR=$project_root/tests/qa_provider_oapi
cd $BUILD_DIR

python3 --version || (echo "We need 'python3' intalled to run integration tests"; exit 1)
pip --version || (echo "We need 'pip' intalled to run integration tests"; exit 1)

pip install -r requirements.txt

pytest -k TF-108 -v ./test_provider_oapi.py
