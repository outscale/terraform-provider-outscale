resource "outscale_vm" "my_vm" {
  image_id           = var.image_id
  vm_type            = var.vm_type
  keypair_name       = outscale_keypair.my_keypair.keypair_name
  security_group_ids = [outscale_security_group.my_sg.security_group_id]
  subnet_id = outscale_subnet.my_subnet.subnet_id

  tags {
    key   = "osc.fcu.eip.auto-attach"
    value = outscale_public_ip.my_public_ip.public_ip
  }
}

resource "local_file" "connect_script" {
  filename        = "${path.module}/connect.sh"
  file_permission = "0770"
  content         = <<EOF
  #!/bin/bash
  ssh -l outscale -o IdentitiesOnly=yes -i ${local_file.my_key.filename} ${outscale_public_ip.my_public_ip.public_ip}
  EOF
}
