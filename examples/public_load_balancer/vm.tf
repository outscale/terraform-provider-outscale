resource "outscale_vm" "my_vms" {
  count                    = var.vm_count
  image_id                 = var.image_id
  vm_type                  = var.vm_type
  placement_subregion_name = "${var.region}a"
  keypair_name             = outscale_keypair.my_keypair.keypair_name
  security_group_ids       = [outscale_security_group.my_sg.security_group_id]

  # quick installaion of a web server for load balancer demo purpose
  provisioner "file" {
    source      = "vm_startup.sh"
    destination = "/tmp/vm_startup.sh"
    connection {
      type        = "ssh"
      user        = "outscale"
      host        = self.public_ip
      private_key = tls_private_key.my_key.private_key_pem
    }
  }

  provisioner "remote-exec" {
    inline = [
      "sudo bash /tmp/vm_startup.sh ${count.index}"
    ]
    connection {
      type        = "ssh"
      user        = "outscale"
      host        = self.public_ip
      private_key = tls_private_key.my_key.private_key_pem
    }
  }
}

resource "local_file" "connect_script" {
  count           = length(outscale_vm.my_vms)
  filename        = "${path.module}/connect_${count.index}.sh"
  file_permission = "0770"
  content         = <<EOF
  #!/bin/bash
  ssh -l outscale -o IdentitiesOnly=yes -i ${local_file.my_key.filename} ${outscale_vm.my_vms[count.index].public_ip}
  EOF
}

