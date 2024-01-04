#!/bin/bash

set -e -x

project_dir=$(cd "$(dirname $0)" && pwd)
project_root=$(cd $project_dir/.. && pwd)
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
           openssl req -x509 -sha256 -nodes -newkey rsa:4096 -keyout test-cert.key -days 2 -out test-cert.pem -subj /CN=domain.com
EOF
    read   = <<-EOF
           echo "{\"filename\":  \"test-cert.pem\"}"
EOF
    delete = ""
  }
  working_directory = path.module
}
' > "outscale/$tf_file"

if [ ! -e "outscale/$tf_file" ]; then
    echo " $tf_file doesn't existe"
    exit 1
fi

cd $project_root/outscale && terraform init || exit 1
cd $project_root/outscale && terraform apply -auto-approve || exit 1

exit 0
