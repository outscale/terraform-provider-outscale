resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF148"
}

resource "outscale_security_group" "sgUsd" {
    description         = "test vms"
    security_group_name = "test-sgUsd"
}

resource "outscale_vm" "outscale_vm1" {
 image_id = var.image_id
 vm_type = var.vm_type
 keypair_name = outscale_keypair.my_keypair.keypair_name
 security_group_names = [outscale_security_group.sgUsd.security_group_name]
 user_data = base64encode(<<EOF
#cloud-config
cloud_config_modules:
- runcmd

runcmd:
- touch /tmp/qa-valid-terraform-user-data-cloud-init
- echo "blabla" >> /dev/ttyS0
EOF
)
tags {
  key = "name"
  value = "Terraform_VM12"
 }
}

