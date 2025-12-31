resource "outscale_net" "outscale_net" {
  ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet01" {
  subregion_name = "${var.region}a"
  ip_range       = "10.0.0.0/24"
  net_id         = outscale_net.outscale_net.net_id
}

resource "outscale_subnet" "subnet02" {
  subregion_name = "${var.region}b"
  ip_range       = "10.0.2.0/24"
  net_id         = outscale_net.outscale_net.net_id
}

resource "outscale_security_group" "security_group01" {
  description         = "Terraform security group test"
  security_group_name = "test-sg-${random_string.suffix[0].result}"
  net_id              = outscale_net.outscale_net.net_id
}

resource "outscale_vm" "vm01" {
  image_id  = var.image_id
  vm_type   = var.vm_type
  subnet_id = outscale_subnet.subnet01.subnet_id
}


resource "outscale_nic" "outscale_nic" {
  subnet_id          = outscale_subnet.subnet01.subnet_id
  description        = "TF-109"
  security_group_ids = [outscale_security_group.security_group01.security_group_id]
  private_ips {
    is_primary = true
    private_ip = "10.0.0.45"
  }

  private_ips {
    is_primary = false
    private_ip = "10.0.0.46"
  }
  tags {
    key   = "Key:"
    value = ":value-tags"
  }
  tags {
    key   = "Key-2"
    value = "value-tags-2"
  }
}

resource "outscale_nic" "outscale_nic_2" {
  subnet_id          = outscale_subnet.subnet01.subnet_id
  security_group_ids = [outscale_security_group.security_group01.security_group_id]
  private_ips {
    is_primary = true
    private_ip = "10.0.0.41"
  }
  private_ips {
    is_primary = false
    private_ip = "10.0.0.42"
  }
  tags {
    key   = "Name"
    value = "Nic-2"
  }
}

resource "outscale_nic" "outscale_nic_3" {
  subnet_id          = outscale_subnet.subnet02.subnet_id
  security_group_ids = [outscale_security_group.security_group01.security_group_id]
  private_ips {
    is_primary = true
    private_ip = "10.0.2.21"
  }
  tags {
    key   = "Key:"
    value = ":value-tags"
  }
  tags {
    key   = "Key-2"
    value = "value-tags-2"
  }
}

resource "outscale_nic_link" "nic_link01" {
  device_number = "1"
  vm_id         = outscale_vm.vm01.vm_id
  nic_id        = outscale_nic.outscale_nic.nic_id
}

resource "outscale_nic_link" "nic_link02" {
  device_number = "2"
  vm_id         = outscale_vm.vm01.vm_id
  nic_id        = outscale_nic.outscale_nic_2.nic_id
}

data "outscale_nic" "nic-0" {
  filter {
    name   = "nic_ids"
    values = [outscale_nic.outscale_nic.nic_id]
  }
}

data "outscale_nic" "nic-1" {
  filter {
    name   = "descriptions"
    values = [outscale_nic.outscale_nic.description]
  }
  filter {
    name   = "states"
    values = [outscale_nic.outscale_nic.state]
  }
  filter {
    name   = "private_ips_primary_ip"
    values = [1]
  }
  filter {
    name   = "private_ips_private_ips"
    values = ["10.0.0.45"]
  }
}

data "outscale_nic" "nic-2-main" {
  filter {
    name   = "link_nic_vm_ids"
    values = [outscale_nic_link.nic_link01.vm_id]
  }
  filter {
    name   = "link_nic_device_numbers"
    values = [0]
  }
}

data "outscale_nic" "nic-4" {
  filter {
    name   = "private_ips_private_ips"
    values = ["10.0.0.42"]
  }
  filter {
    name   = "tag_keys"
    values = ["Name"]
  }
  filter {
    name   = "tag_values"
    values = ["Nic-2"]
  }
  depends_on = [outscale_nic.outscale_nic, outscale_nic.outscale_nic_2, outscale_nic.outscale_nic_3]
}

data "outscale_nic" "nic-5" {
  filter {
    name   = "security_group_ids"
    values = [outscale_security_group.security_group01.security_group_id]
  }
  filter {
    name   = "tag_keys"
    values = ["Key-2"]
  }
  filter {
    name   = "tag_values"
    values = ["value-tags-2"]
  }
  filter {
    name   = "private_ips_private_ips"
    values = ["10.0.0.46"]
  }

  depends_on = [outscale_nic.outscale_nic, outscale_nic.outscale_nic_2, outscale_nic.outscale_nic_3]
}

data "outscale_nic" "nic-6" {
  filter {
    name   = "security_group_names"
    values = [outscale_security_group.security_group01.security_group_name]
  }
  filter {
    name   = "private_ips_private_ips"
    values = ["10.0.2.21"]
  }
  filter {
    name   = "private_ips_primary_ip"
    values = ["true"]
  }

  depends_on = [outscale_nic.outscale_nic, outscale_nic.outscale_nic_2, outscale_nic.outscale_nic_3]
}
