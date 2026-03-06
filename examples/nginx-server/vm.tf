# Create the virtual machine and initialize it with cloud-init.
resource "outscale_vm" "my-vm" {
  vm_type            = var.vm_type
  image_id           = var.image_id
  # keypair_name_wo is the write-only version of keypair_name.
  # It avoids storing the value in Terraform state.
  keypair_name_wo    = outscale_keypair.my-keypair.keypair_name
  security_group_ids = [outscale_security_group.my-security-group.security_group_id]

  # Bootstrap the instance at first boot.
  user_data = local.user_data

  # Resize the root volume.
  block_device_mappings {
    device_name = "/dev/sda1"

    bsu {
      volume_size           = var.root_volume_size
      volume_type           = "gp2"
      delete_on_vm_deletion = true
    }
  }

  dynamic "tags" {
    for_each = local.common_tags

    content {
      key   = tags.key
      value = tags.value
    }
  }
}