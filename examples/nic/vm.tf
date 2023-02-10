resource "outscale_vm" "my_vm" {
  image_id                 = var.image_id
  vm_type                  = var.vm_type
  keypair_name             = outscale_keypair.my_keypair.keypair_name
  
  nics {
    nic_id = outscale_nic.eth0.nic_id
    device_number = 0
  }

  nics {
    nic_id = outscale_nic.eth1.nic_id
    device_number = 1
  }

  # resized bootdisk volume
  block_device_mappings {
    device_name = "/dev/sda1"
    bsu {
      volume_size           = "100"
      volume_type           = "gp2"
      delete_on_vm_deletion = "true"
    }
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
