resource "outscale_net" "net01" {
  ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet01" {
  net_id         = outscale_net.net01.net_id
  ip_range       = "10.0.0.0/24"
  subregion_name = "eu-west-2a"
}

resource "outscale_nic" "nic01" {
  subnet_id = outscale_subnet.subnet01.subnet_id
}

resource "outscale_subnet" "subnet02" {
  net_id         = outscale_net.net01.net_id
  ip_range       = "10.0.1.0/24"
  subregion_name = "eu-west-2a"
}

resource "outscale_nic" "nic02" {
  subnet_id = outscale_subnet.subnet02.subnet_id
}

resource "outscale_nic" "nic03" {
  subnet_id = outscale_subnet.subnet02.subnet_id
}

resource "outscale_keypair" "keypair01" {
  keypair_name = "terraform-keypair-for-vm"
}

resource "outscale_vm" "vm01" {
  image_id     = var.image_id
  vm_type      = "tinav7.c1r1p2"
  keypair_name = outscale_keypair.keypair01.keypair_name
  primary_nic {
    nic_id        = outscale_nic.nic01.nic_id
    device_number = "0"
  }
}

resource "outscale_nic_link" "nic_link02" {
  device_number = "2"
  vm_id         = outscale_vm.vm01.vm_id
  nic_id        = outscale_nic.nic03.nic_id
}
