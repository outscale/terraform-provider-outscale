resource "outscale_vm" "my_vm" {
  image_id                 = var.image_id
  vm_type                  = var.vm_type
  keypair_name             = outscale_keypair.my_keypair.keypair_name
  security_group_ids       = [outscale_security_group.my_sg.security_group_id]
  placement_subregion_name = "${var.region}a"
  placement_tenancy        = "default"
  get_admin_password       = true

  # resized bootdisk volume
  block_device_mappings {
    device_name = "/dev/sda1"
    bsu {
      volume_size           = "100"
      volume_type           = "gp2"
      delete_on_vm_deletion = "true"
    }
  }

  tags {
    key   = "osc.fcu.eip.auto-attach"
    value = outscale_public_ip.my_public_ip.public_ip
  }
}

resource "local_file" "password_txt" {
  filename = "${path.module}/password.txt"
  content  = rsadecrypt(outscale_vm.my_vm.admin_password, tls_private_key.my_key.private_key_pem)
}
