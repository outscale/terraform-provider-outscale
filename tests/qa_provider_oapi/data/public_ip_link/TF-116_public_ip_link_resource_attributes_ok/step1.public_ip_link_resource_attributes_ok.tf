resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF116"
}
resource "outscale_public_ip" "outscale_public_ip" {
 tags {
      key = "name"
      value = "public_ip"
      }
}

resource "outscale_vm" "outscale_vm" {
    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = outscale_keypair.my_keypair.keypair_name
}

resource "outscale_public_ip_link" "outscale_public_ip_link" {
    vm_id             = outscale_vm.outscale_vm.vm_id
    public_ip          = outscale_public_ip.outscale_public_ip.public_ip
}
