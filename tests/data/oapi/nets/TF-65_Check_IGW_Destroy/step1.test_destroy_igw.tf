resource "outscale_internet_service" "nomad" {
depends_on = [outscale_net.nomad]
}
resource "outscale_net" "nomad" {
  ip_range = "10.0.0.0/16"
  tags {
    key   = "name"
    value = "TF65-network"
  }
}
resource "outscale_internet_service_link" "nomad" {
  internet_service_id = outscale_internet_service.nomad.internet_service_id
  net_id              = outscale_net.nomad.net_id
}
resource "outscale_subnet" "bastion" {
  subregion_name = "${var.region}a"
  ip_range       = "10.0.0.0/24"
  net_id         = outscale_net.nomad.net_id
  tags {
    key   = "name"
    value = "TF65-bastion"
  }
}
resource "outscale_subnet" "adm" {
  subregion_name = "${var.region}a"
  ip_range       = "10.0.1.0/24"
  net_id         = outscale_net.nomad.net_id
  tags {
    key   = "name"
    value = "TF65-adm"
  }
}
resource "outscale_public_ip" "nat" {
     tags {
        key = "name"
        value = "test1"
      }
}
resource "outscale_nat_service" "adm" {
  depends_on = [outscale_internet_service.nomad]
  subnet_id    = outscale_subnet.bastion.subnet_id
  public_ip_id = outscale_public_ip.nat.public_ip_id
}

resource "outscale_route_table" "public" {
  net_id = outscale_net.nomad.net_id
}
resource "outscale_route" "internet" {
  destination_ip_range = "0.0.0.0/0"
  gateway_id           = outscale_internet_service.nomad.internet_service_id
  route_table_id       = outscale_route_table.public.route_table_id
}
resource "outscale_route_table_link" "igw" {
  subnet_id      = outscale_subnet.bastion.subnet_id
  route_table_id = outscale_route_table.public.id
}
resource "outscale_route_table" "nat" {
  net_id = outscale_net.nomad.net_id
}
resource "outscale_route" "nat_internet" {
  destination_ip_range = "0.0.0.0/0"
  nat_service_id = outscale_nat_service.adm.nat_service_id
  route_table_id       = outscale_route_table.nat.route_table_id
depends_on = [outscale_route.internet]
}
resource "outscale_route_table_link" "nat" {
  subnet_id      = outscale_subnet.adm.subnet_id
  route_table_id = outscale_route_table.nat.id
}
resource "outscale_vm" "consul_server_1" {
 count = 3
  security_group_ids = [outscale_security_group.nomad-sg2.security_group_id]
 image_id     = var.image_id
  vm_type      = var.vm_type
  subnet_id = outscale_subnet.bastion.subnet_id
  tags {
    key   = "name"
    value = "consul-server-1"
  }
}
resource "outscale_vm" "vm2" {
 count = 3
  security_group_ids = [outscale_security_group.nomad-sg1.security_group_id]
  image_id     = var.image_id
  vm_type      = var.vm_type
  subnet_id = outscale_subnet.adm.subnet_id
  tags {
    key   = "name"
    value = "consul-server-1"
  }
}
resource "outscale_security_group" "nomad-sg1" {
                description         = "sg for terraform tests"
                security_group_name = "test-sg-${random_string.suffix[0].result}"
                net_id              = outscale_net.nomad.net_id
        }
resource "outscale_security_group" "nomad-sg2" {
                description         = "sg for terraform tests"
                security_group_name = "test-sg-${random_string.suffix[1].result}"
                net_id              = outscale_net.nomad.net_id
        }

resource "outscale_security_group_rule" "rule1-sg1" {
    flow              = "Inbound"
    security_group_id = outscale_security_group.nomad-sg1.id
    from_port_range   = 0
    to_port_range     = 0
    ip_protocol       = "tcp"
    ip_range          = "192.168.0.1/32"
depends_on =[outscale_security_group.nomad-sg1]
}

resource "outscale_security_group_rule" "rule1-sg2" {
    flow              = "Inbound"
    security_group_id = outscale_security_group.nomad-sg2.id
    from_port_range   = 0
    to_port_range     = 0
    ip_protocol       = "tcp"
    ip_range          = "192.168.0.1/32"
depends_on =[outscale_security_group.nomad-sg2]
}

resource "outscale_security_group_rule" "rule2-sg1" {
    flow              = "Inbound"
    security_group_id = outscale_security_group.nomad-sg1.id
    from_port_range   = 22
    to_port_range     = 22
    ip_protocol       = "tcp"
    ip_range          = "192.168.0.1/32"
depends_on =[outscale_security_group_rule.rule1-sg1]
}

resource "outscale_security_group_rule" "rule2-sg2" {
    flow              = "Inbound"
    security_group_id = outscale_security_group.nomad-sg2.id
    from_port_range   = 80
    to_port_range     = 80
    ip_protocol       = "tcp"
    ip_range          = "192.168.0.1/32"
depends_on =[outscale_security_group_rule.rule1-sg2]
}
resource "outscale_security_group_rule" "rule3-sg1" {
    flow              = "Inbound"
    security_group_id = outscale_security_group.nomad-sg1.id
    from_port_range   = 1024
    to_port_range     = 1024
    ip_protocol       = "tcp"
    ip_range          = "192.168.0.1/32"
depends_on =[outscale_security_group_rule.rule2-sg1]
}

resource "outscale_security_group_rule" "rule3-sg2" {
    flow              = "Inbound"
    security_group_id = outscale_security_group.nomad-sg2.id
    from_port_range   = 1024
    to_port_range     = 1024
    ip_protocol       = "tcp"
    ip_range          = "192.168.0.1/32"
depends_on =[outscale_security_group_rule.rule2-sg2]
}

resource "outscale_public_ip" "EIP" {
     tags {
        key = "name"
        value = "EIP-TF65"
      }
}
resource "outscale_public_ip_link" "eip_link" {
 vm_id             = outscale_vm.consul_server_1[0].vm_id
 public_ip          = outscale_public_ip.EIP.public_ip
}

resource "outscale_public_ip" "EIP2" {
     tags {
        key = "name"
        value = "EIP-TF65-2"
      }
}
resource "outscale_public_ip_link" "eip_link2" {
 vm_id             = outscale_vm.consul_server_1[1].vm_id
 public_ip          = outscale_public_ip.EIP2.public_ip
}
