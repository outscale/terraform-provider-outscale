resource "outscale_public_ip" "my_public_ip" {
}

resource "outscale_public_ip_link" "my_public_ip_link" {
  vm_id     = outscale_vm.my_vm.vm_id
  public_ip = outscale_public_ip.my_public_ip.public_ip
}

resource "local_file" "ip_txt" {
  filename = "${path.module}/ip.txt"
  content  = outscale_public_ip.my_public_ip.public_ip
}
