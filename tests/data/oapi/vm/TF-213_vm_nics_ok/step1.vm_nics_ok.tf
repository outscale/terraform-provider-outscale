resource "outscale_net" "net01" {
  ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet01" {
  net_id         = outscale_net.net01.net_id
  ip_range       = "10.0.0.0/24"
  subregion_name = "eu-west-2a"
}

resource "outscale_subnet" "subnet02" {
  net_id         = outscale_net.net01.net_id
  ip_range       = "10.0.1.0/24"
  subregion_name = "eu-west-2a"
}

resource "outscale_keypair" "keypair01" {
  keypair_name = "terraform-keypair-for-vm"
}

resource "outscale_vm" "vm01" {
  image_id     = var.image_id
  vm_type      = "tinav7.c1r1p2"
  keypair_name = outscale_keypair.keypair01.keypair_name
  nics {
    delete_on_vm_deletion = true
    subnet_id     = outscale_subnet.subnet01.id
    device_number = "0"
  }
  nics {
    delete_on_vm_deletion = true
    subnet_id     = outscale_subnet.subnet02.id
    device_number = "1"
  }
  nics {
    delete_on_vm_deletion = true
    subnet_id     = outscale_subnet.subnet02.id
    device_number = "2"
  }
}
