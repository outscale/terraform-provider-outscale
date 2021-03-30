resource "outscale_public_ip" "outscale_public_ip" {
 tags {
      key = "name"
      value = "public_ip"
      }
}

resource "outscale_vm" "outscale_vm" {
    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = var.keypair_name
    security_group_ids = [var.security_group_id]
}

resource "outscale_public_ip_link" "outscale_public_ip_link" {
    vm_id             = outscale_vm.outscale_vm.vm_id
    #vm_id              = outscale_vm.outscale_vm.id                        # test purposes only
    public_ip          = outscale_public_ip.outscale_public_ip.public_ip
}
