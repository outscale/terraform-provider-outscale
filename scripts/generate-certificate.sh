#!/bin/bash

set -e

project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)

if [ ! -d "$project_root/tests/certs" ]; then
    mkdir $project_root/tests/certs
fi
build_dir=$(cd $project_root/tests/certs && pwd)
tf_file="gen-cert-test.tf"

cd $build_dir
echo '
terraform {
  required_providers {
    shell = {
      source  = "scottwinkler/shell"
      version = ">=1.7.10"
    }
  }
}

resource "shell_script" "ca_gen" {
  lifecycle_commands {
    create = <<-EOF
           openssl req -x509 -sha256 -nodes -newkey rsa:4096 -keyout certificate.key -days 2 -out certificate.pem -subj /CN=domain.com
EOF
    read   = <<-EOF
           echo "{\"filename\":  \"certificate.pem\"}"
EOF
    delete = ""
  }
  working_directory = path.module
}
' > "$build_dir/$tf_file"

if [ ! -e "$build_dir/$tf_file" ]; then
    echo "$tf_file doesn't exist"
    exit 1
fi

echo "Generating certificates in $build_dir"
terraform init || exit 1
terraform apply -auto-approve || exit 1

oapi_testdata="$project_root/internal/services/oapi/testdata"
if [ ! -d "$oapi_testdata" ]; then
    mkdir $oapi_testdata
fi
cp certificate.pem certificate.key $oapi_testdata/

cd $project_root

exit 0
