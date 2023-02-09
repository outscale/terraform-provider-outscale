resource "outscale_keypair" "keypair-TF66" {
  count        = 2
  keypair_name = "keyname_TF66-${count.index}"
}

resource "outscale_security_group" "security_group_TF66" {
  count               = 2
  description         = "test-terraform-TF66"
  security_group_name = "terraform-sg-TF66-${count.index}"
}

resource "outscale_vm" "outscale_vm-TF66" {
  block_device_mappings {
    device_name                  = "/dev/sdb"
      bsu {
        volume_size              = 20
        volume_type              = "standard"
      }
    }
  image_id                       = var.image_id
  vm_type                        = "tinav4.c2r2p3"
  performance                    = "medium"
  deletion_protection            = false
  vm_initiated_shutdown_behavior = "restart"
  security_group_ids             = [outscale_security_group.security_group_TF66[0].security_group_id,outscale_security_group.security_group_TF66[1].security_group_id]
  keypair_name                   = outscale_keypair.keypair-TF66[1].keypair_name
}
