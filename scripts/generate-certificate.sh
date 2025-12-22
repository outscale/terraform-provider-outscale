#!/bin/bash

set -e

project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
if [ ! -d "$project_root/tests/data/cert_example" ]; then
    mkdir $project_root/tests/data/cert_example
fi
build_dir=$(cd $project_root/tests/data/cert_example && pwd)
tf_file="gen-cert-test.tf"

cd $project_root
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
' | tee "$project_root/outscale/$tf_file"  "$build_dir/$tf_file"

if [ ! -e "$build_dir/$tf_file" ] && [ ! -e "$project_root/outscale/$tf_file" ]; then
    echo " $tf_file doesn't existe"
    exit 1
fi

cd outscale/
terraform init || exit 1
terraform apply -auto-approve || exit 1

cd $build_dir
terraform init || exit 1
terraform apply -auto-approve || exit 1
cd $project_root

exit 0
