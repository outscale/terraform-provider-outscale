resource "outscale_keypair" "my_keypair" {
 keypair_name = "test-keypair-${random_string.suffix[0].result}"
}
resource "outscale_public_ip" "outscale_public_ip" {
 tags {
      key = "name"
      value = "public_ip"
      }
}

resource "outscale_security_group" "sgPub" {
   description         = "sg for terraform tests"
   security_group_name = "test-sg-${random_string.suffix[0].result}"
}


resource "outscale_vm" "outscale_vm" {
    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = outscale_keypair.my_keypair.keypair_name
    security_group_ids = [outscale_security_group.sgPub.security_group_id]
}

resource "outscale_public_ip_link" "outscale_public_ip_link" {
    vm_id             = outscale_vm.outscale_vm.vm_id
    public_ip          = outscale_public_ip.outscale_public_ip.public_ip
}
