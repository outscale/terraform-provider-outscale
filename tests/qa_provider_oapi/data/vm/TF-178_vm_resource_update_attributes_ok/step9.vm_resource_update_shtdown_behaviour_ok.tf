resource "outscale_keypair" "keypair-TF178" {
  count        = 2
  keypair_name = "keyname_TF178-${count.index}"
}

resource "outscale_security_group" "security_group_TF178" {
  count               = 2
  description         = "test-terraform-TF178"
  security_group_name = "terraform-sg-TF178-${count.index}"
}

resource "outscale_vm" "outscale_vm-TF178" {
  block_device_mappings {
    device_name                  = "/dev/sdb"
      bsu {
        volume_size              = 20
        volume_type              = "standard"
        delete_on_vm_deletion    = true
      }
    }
  image_id                       = var.image_id
  vm_type                        = "tinav5.c3r3"
  performance                    = "high"
  deletion_protection            = false
  vm_initiated_shutdown_behavior = "restart"
  security_group_ids             = [outscale_security_group.security_group_TF178[0].security_group_id,outscale_security_group.security_group_TF178[1].security_group_id]
  keypair_name                   = outscale_keypair.keypair-TF178[1].keypair_name
  user_data                      = "LS0tLS1CRUdJTiBPVVRTQ0FMRSBTRUNUSU9OLS0tLS0KCnByaXZhdGVfb25seT10cnVlCgotLS0tLUVORCBPVVRTQ0FMRSBTRUNUSU9OLS0tLS0="
}
