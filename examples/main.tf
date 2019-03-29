# resource "outscale_lin" "outscale_lin1" {
#   cidr_block = "10.0.0.0/16"
# }
# resource "outscale_route_table" "outscale_route_table1" {
#   vpc_id = "${outscale_lin.outscale_lin1.vpc_id}"
# }
# resource "outscale_route_table" "outscale_route_table2" {
#   vpc_id = "${outscale_lin.outscale_lin1.vpc_id}"
# }
# data "outscale_route_tables" "outscale_route_tables" {
#   route_table_id = ["${outscale_route_table.outscale_route_table1.route_table_id}", "${outscale_route_table.outscale_route_table2.route_table_id}"]
# }
# resource "outscale_lin" "outscale_lin2" {
#   count = 1
#   cidr_block = "10.0.0.0/16"
# }
# resource "outscale_subnet" "outscale_subnet1" {
#   count = 1
#   availability_zone = "eu-west-2a"
#   cidr_block        = "10.0.0.0/16"
#   vpc_id            = "${outscale_lin.outscale_lin2.vpc_id}"
# }
# output "outscale_subnet1" {
#   value = "${outscale_subnet.outscale_subnet1.subnet_id}"
# }
# resource "outscale_lin" "outscale_lin" {
#   cidr_block = "10.0.0.0/24"
# }
# resource "outscale_route_table" "outscale_route_table" {
#   vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
# }
# resource "outscale_lin_internet_gateway" "outscale_lin_internet_gateway" {}
# resource "outscale_lin_internet_gateway_link" "outscale_lin_internet_gateway_link" {
#   internet_gateway_id = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.internet_gateway_id}"
#   vpc_id              = "${outscale_lin.outscale_lin.vpc_id}"
# }
# resource "outscale_route" "outscale_route" {
#   gateway_id             = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.internet_gateway_id}"
#   destination_cidr_block = "10.0.0.0/16"
#   route_table_id         = "${outscale_route_table.outscale_route_table.route_table_id}"
# }
# resource "outscale_keypair_importation" "outscale_keypair_importation" {
#   key_name            = "keyname_test"
#   public_key_material = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIj41qgxMCP/iteneqXSIFZBp5vizPvaoIR3Um9xK7PGoW8giupGn+EPuxIA4cDM4vzOqOkiMPhz5XK0whEjkVzTo4+S0puvDZuwIsdiW9mxhJc7tgBNL0cYlWSYVkz4G/fslNfRPW5mYAM49f4fhtxPb5ok4Q2Lg9dPKVHO/Bgeu5woMc7RY0p1ej6D4CKFE6lymSDJpW0YHX/wqE9+cfEauh7xZcG0q9t2ta6F6fmX0agvpFyZo8aFbXeUBr7osSCJNgvavWbM/06niWrOvYX2xwWdhXmXSrbX8ZbabVohBK41 phodgson@thoughtworks.com"
# }
# resource "outscale_lin" "outscale_lin" {
#   cidr_block = "10.0.0.0/16"
# }
# resource "outscale_subnet" "outscale_subnet" {
#   vpc_id     = "${outscale_lin.outscale_lin.vpc_id}"
#   cidr_block = "10.0.0.0/18"
# }
# resource "outscale_public_ip" "outscale_public_ip" {
#   #domain = "Standard" # BUG doc API
#   domain = ""
# }
# resource "outscale_nat_service" "outscale_nat_service" {
#   depends_on    = ["outscale_route.outscale_route"]
#   subnet_id     = "${outscale_subnet.outscale_subnet.subnet_id}"
#   allocation_id = "${outscale_public_ip.outscale_public_ip.allocation_id}"
# }
# resource "outscale_route_table" "outscale_route_table" {
#   vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
# }
# resource "outscale_route" "outscale_route" {
#   destination_cidr_block = "0.0.0.0/0"
#   gateway_id             = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.id}"
#   route_table_id         = "${outscale_route_table.outscale_route_table.id}"
# }
# resource "outscale_route_table_link" "outscale_route_table_link" {
#   subnet_id      = "${outscale_subnet.outscale_subnet.subnet_id}"
#   route_table_id = "${outscale_route_table.outscale_route_table.id}"
# }
# resource "outscale_lin_internet_gateway" "outscale_lin_internet_gateway" {}
# resource "outscale_lin_internet_gateway_link" "outscale_lin_internet_gateway_link" {
#   vpc_id              = "${outscale_lin.outscale_lin.vpc_id}"
#   internet_gateway_id = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.id}"
# }
# resource "outscale_client_endpoint" "outscale_client_endpoint" {
#   bgp_asn    = "3"
#   ip_address = "171.33.74.122"
#   type       = "ipsec.1"
# }
# resource "outscale_dhcp_option" "outscale_dhcp_option" {}
# resource "outscale_lin" "outscale_lin" {
#   cidr_block = "10.0.0.0/16"
# }
# resource "outscale_dhcp_option_link" "outscale_dhcp_option_link" {
#   dhcp_options_id = "${outscale_dhcp_option.outscale_dhcp_option.dhcp_options_id}"
#   vpc_id          = "${outscale_lin.outscale_lin.vpc_id}"
# }
# resource "outscale_keypair" "outscale_keypair" {
#   count = 1
#   key_name = "keyname_test_123"
# }
# resource "outscale_volume" "outscale_volume" {
#   availability_zone = "eu-west-2a"
#   size              = 40
# }
# resource "outscale_vm" "outscale_vm" {
#   image_id      = "ami-880caa66"
#   instance_type = "c4.large"
#   # key_name       = "integ_sut_keypair"
#   # security_group = ["sg-c73d3b6b"]
# }
# resource "outscale_volumes_link" "outscale_volumes_link" {
#   device      = "/dev/sdb"
#   volume_id   = "${outscale_volume.outscale_volume.id}"
#   instance_id = "${outscale_vm.outscale_vm.id}"
# }
# resource "outscale_lin" "outscale_lin" {
#   cidr_block = "10.0.0.0/16"
# }
# resource "outscale_subnet" "outscale_subnet" {
#   vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
#   cidr_block = "10.0.0.0/18"
# }
# resource "outscale_public_ip" "outscale_public_ip" {
#   #domain               = "Standard"       # BUG doc API
#   domain = ""
# }
# resource "outscale_nat_service" "outscale_nat_service" {
#   depends_on = ["outscale_route.outscale_route"]
#   subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
#   allocation_id = "${outscale_public_ip.outscale_public_ip.allocation_id}"
# }
# resource "outscale_route_table" "outscale_route_table" {
#   vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
# }
# resource "outscale_route" "outscale_route" {
#   destination_cidr_block = "0.0.0.0/0"
#   gateway_id = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.id}"
#   route_table_id = "${outscale_route_table.outscale_route_table.id}"
# }
# resource "outscale_route_table_link" "outscale_route_table_link" {
#   subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
#   route_table_id = "${outscale_route_table.outscale_route_table.id}"
# }
# resource "outscale_lin_internet_gateway" "outscale_lin_internet_gateway" {}
# resource "outscale_lin_internet_gateway_link" "outscale_lin_internet_gateway_link" {
#   vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
#   internet_gateway_id = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.id}"
# }
# data "outscale_nat_service" "outscale_nat_service" {
#   nat_gateway_id = "${outscale_nat_service.outscale_nat_service.nat_gateway_id}"
# }
# resource "outscale_volume" "outscale_volume" {
#   availability_zone = "eu-west-2a"
#   size = 40
# }
# data "outscale_volume" "outscale_volume" {
#   volume_id = "${outscale_volume.outscale_volume.volume_id}"
# }
# resource "outscale_client_endpoint" "outscale_client_endpoint" {
#   bgp_asn = "3"
#   ip_address = "171.33.74.122"
#   type = "ipsec.1"
# }
# data "outscale_client_endpoints" "outscale_client_endpoints" {
#   depends_on = ["outscale_client_endpoint.outscale_client_endpoint"]
#   filter {
#     name = "ip-address"
#     values = ["171.33.74.122"]
#   }
#   /* filter {
#         name = "bgp-asn"
#         values = ["3"] }
#     filter {
#         name = "type"
#         values = ["ipsec.1"] }
#     */
# }
# resource "outscale_client_endpoint" "outscale_client_endpoint" {
#   bgp_asn    = "3"
#   ip_address = "171.33.74.122"
#   type       = "ipsec.1"
#   }
# resource "outscale_lin" "outscale_lin" {
#   count = 1
#   cidr_block = "10.0.0.0/16"
# }
# resource "outscale_vpn_gateway" "outscale_vpn_gateway" {
#   type = "ipsec.1"
# }
# resource "outscale_vpn_gateway_link" "test" {
#   vpc_id         = "${outscale_lin.outscale_lin.id}"
#   vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.id}"
# }
# resource "outscale_route_table" "outscale_route_table" {
#   count = 1
#   vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
# }
# resource "outscale_vpn_gateway_route_propagation" "foo" {
#   gateway_id     = "${outscale_vpn_gateway.outscale_vpn_gateway.vpn_gateway_id}"
#   route_table_id = "${outscale_route_table.outscale_route_table.route_table_id}"
# }
# resource "outscale_snapshot_export_task" "outscale_snapshot_export_task" {
#   count = 1
#   export_to_osu {
#     disk_image_format = "raw"
#     osu_bucket        = "test"
#   }
#   snapshot_id = "snap-5bcc0764"
# }
# resource "outscale_lin" "outscale_lin" {
#   cidr_block = "10.0.0.0/16"
# }
# resource "outscale_subnet" "outscale_subnet" {
#   vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
#   cidr_block = "10.0.0.0/18"
# }
# resource "outscale_public_ip" "outscale_public_ip" {
#   #domain               = "Standard"       # BUG doc API
#   domain = ""
# }
# resource "outscale_nat_service" "outscale_nat_service" {
#   depends_on = ["outscale_route.outscale_route"]
#   subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
#   allocation_id = "${outscale_public_ip.outscale_public_ip.allocation_id}"
# }
# resource "outscale_route_table" "outscale_route_table" {
#   vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
# }
# resource "outscale_route" "outscale_route" {
#   destination_cidr_block = "0.0.0.0/0"
#   gateway_id = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.id}"
#   route_table_id = "${outscale_route_table.outscale_route_table.id}"
# }
# resource "outscale_route_table_link" "outscale_route_table_link" {
#   subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
#   route_table_id = "${outscale_route_table.outscale_route_table.id}"
# }
# resource "outscale_lin_internet_gateway" "outscale_lin_internet_gateway" {}
# resource "outscale_lin_internet_gateway_link" "outscale_lin_internet_gateway_link" {
#   vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
#   internet_gateway_id = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.id}"
# }
# data "outscale_nat_service" "outscale_nat_service" {
#   nat_gateway_id = "${outscale_nat_service.outscale_nat_service.nat_gateway_id}"
# }
# resource "outscale_vm" "outscale_vm" {
#   image_id                = "ami-880caa66"
#   instance_type           = "c4.large"
#   disable_api_termination = false
# }
# data "outscale_vms_state" "outscale_vms_state" {
#   instance_id = ["${outscale_vm.outscale_vm.id}"]
# }
# data "outscale_vm_state" "outscale_vm_state" {
#   instance_id = ["${outscale_vm.outscale_vm.id}"]
# }
# data "outscale_vpn_connection" "outscale_vpn_connection1" {
#   vpn_connection_id = "${outscale_vpn_connection.outscale_vpn_connection.id}"
# }
# data "outscale_vpn_connections" "outscale_vpn_connections" {
#   vpn_connection_id = ["${outscale_vpn_connection.outscale_vpn_connection.id}", "${outscale_vpn_connection.outscale_vpn_connection2.id}"]
# }
# resource "outscale_vpn_connection" "outscale_vpn_connection" {
#   customer_gateway_id = "${outscale_client_endpoint.outscale_client_endpoint.id}"
#   vpn_gateway_id      = "${outscale_vpn_gateway.outscale_vpn_gateway.id}"
#   type                = "ipsec.1"
# }
# resource "outscale_vpn_gateway" "outscale_vpn_gateway" {
#   type = "ipsec.1"
# }
# resource "outscale_client_endpoint" "outscale_client_endpoint" {
#   bgp_asn    = "3"
#   ip_address = "171.33.74.125"
#   type       = "ipsec.1"
# }
# resource "outscale_vpn_gateway_link" "outscale_vpn_gateway_link" {
#   vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
#   #vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.vpn_gateway_id}"
#   vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.id}"
# }
# resource "outscale_lin" "outscale_lin" {
#   cidr_block = "10.0.0.0/16"
# }
# resource "outscale_vpn_connection" "outscale_vpn_connection2" {
#   customer_gateway_id = "${outscale_client_endpoint.outscale_client_endpoint2.id}"
#   vpn_gateway_id      = "${outscale_vpn_gateway.outscale_vpn_gateway2.id}"
#   type                = "ipsec.1"
# }
# resource "outscale_vpn_gateway" "outscale_vpn_gateway2" {
#   type = "ipsec.1"
# }
# resource "outscale_client_endpoint" "outscale_client_endpoint2" {
#   bgp_asn    = "3"
#   ip_address = "171.33.74.126"
#   type       = "ipsec.1"
# }
# resource "outscale_vpn_gateway_link" "outscale_vpn_gateway_link2" {
#   vpc_id = "${outscale_lin.outscale_lin2.vpc_id}"
#   #vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway2.vpn_gateway_id}"
#   vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway2.id}"
# }
# resource "outscale_lin" "outscale_lin2" {
#   cidr_block = "10.0.0.0/16"
# }
# resource "outscale_volume" "test" {
#   availability_zone = "eu-west-2a"
#   size              = 1
# }
# resource "outscale_snapshot" "test" {
#   volume_id = "${outscale_volume.test.id}"
# }
# resource "outscale_load_balancer" "outscale_load_balancer" {
#   count = 1
#   load_balancer_name = "foobar-terraform-elb"
#   availability_zones = ["eu-west-2a"]
#   listeners {
#     instance_port = 1024
#     instance_protocol = "HTTP"
#     load_balancer_port = 25
#     protocol = "HTTP"
#   }
# }
# resource "outscale_load_balancer" "outscale_load_balancer2" {
#   count = 1
#   load_balancer_name = "foobar-terraform-elb2"
#   availability_zones = ["eu-west-2a"]
#   listeners {
#     instance_port = 1024
#     instance_protocol = "HTTP"
#     load_balancer_port = 25
#     protocol = "HTTP"
#   }
# }
# data "outscale_load_balancers" "outscale_load_balancers" {
#   load_balancer_name = ["${outscale_load_balancer.outscale_load_balancer.load_balancer_name}", "${outscale_load_balancer.outscale_load_balancer2.load_balancer_name}"]
# }
# resource "outscale_lin" "outscale_lin" {
#   count = 1
#   cidr_block = "10.0.0.0/16"
# }
# resource "outscale_vpn_gateway" "outscale_vpn_gateway" {
#   type = "ipsec.1"
# }
# resource "outscale_vpn_gateway_link" "test" {
#   vpc_id         = "${outscale_lin.outscale_lin.id}"
#   vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.id}"
# }
# resource "outscale_route_table" "outscale_route_table" {
#   count = 1
#   vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
# }
# resource "outscale_vpn_gateway_route_propagation" "foo" {
#   gateway_id     = "${outscale_vpn_gateway.outscale_vpn_gateway.vpn_gateway_id}"
#   route_table_id = "${outscale_route_table.outscale_route_table.route_table_id}"
# }

# resource "outscale_vm" "basic" {
#   image_id      = "ami-880caa66"
#   instance_type = "t2.micro"
# }

# resource "outscale_image" "foo" {
#   name        = "tf-testing-foo"
#   instance_id = "${outscale_vm.basic.id}"
# }

# resource "outscale_volume" "outscale_volume" {
#   availability_zone = "eu-west-2a"
#   size              = 40
# }

# resource "outscale_snapshot" "outscale_snapshot" {
#   volume_id = "${outscale_volume.outscale_volume.volume_id}"
# }

# resource "outscale_image_register" "outscale_image_register" {
#   name = "registeredImageFromSnapshot"

#   root_device_name = "/dev/sda1"

#   block_device_mapping {
#     ebs {
#       snapshot_id = "${outscale_snapshot.outscale_snapshot.snapshot_id}"
#     }
#   }
# }

# resource "outscale_volume" "example" {
#   availability_zone = "eu-west-2a"
#   volume_type       = "gp2"
#   size              = 40

#   tag {
#     Name = "External Volume"
#   }
# }

# resource "outscale_snapshot" "snapshot" {
#   volume_id = "${outscale_volume.example.id}"
# }

# data "outscale_snapshot" "snapshot" {
#   snapshot_id = "${outscale_snapshot.snapshot.id}"
# }

# resource "outscale_vm" "outscale_instance" {
#   image_id      = "ami-880caa66"
#   instance_type = "c4.large"
#   subnet_id     = "${outscale_subnet.outscale_subnet.subnet_id}"
# }

resource "outscale_lin" "outscale_lin" {
  cidr_block = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
  availability_zone = "eu-west-2a"
  cidr_block        = "10.0.0.0/16"
  vpc_id            = "${outscale_lin.outscale_lin.id}"
}

resource "outscale_nic" "outscale_nic" {
  subnet_id         = "${outscale_subnet.outscale_subnet.subnet_id}"
  security_group_id = ["${outscale_firewall_rules_set.web.id}"]
}

resource "outscale_nic_link" "outscale_nic_link" {
  device_index         = "1"
  instance_id          = "${outscale_vm.basic.id}"
  network_interface_id = "${outscale_nic.outscale_nic.id}"
}

resource "outscale_nic_private_ip" "outscale_nic_private_ip" {
  network_interface_id               = "${outscale_nic.outscale_nic.id}"
  secondary_private_ip_address_count = 5
}

resource "outscale_keypair" "a_key_pair" {
  key_name = "terraform-key-test21"
}

resource "outscale_firewall_rules_set" "web" {
  group_name        = "lin_ucP2_sg_allow_me"
  group_description = "Allow inbound traffic from me"
  vpc_id            = "${outscale_lin.outscale_lin.id}"

  tag {
    Name = "lin_ucP2_sg_allow_me"
  }
}

resource "outscale_inbound_rule" "allow_men2" {
  ip_permissions = {
    ip_protocol = "tcp"
    from_port   = 22
    to_port     = 22
    ip_ranges   = ["10.0.0.0/16"]
  }

  group_id = "${outscale_firewall_rules_set.web.id}"
}

resource "outscale_vm" "basic" {
  image_id       = "ami-880caa66"
  instance_type  = "c4.large"
  subnet_id      = "${outscale_subnet.outscale_subnet.subnet_id}"
  key_name       = "${outscale_keypair.a_key_pair.key_name}"
  security_group = ["${outscale_firewall_rules_set.web.id}"]
}

resource "outscale_lin_internet_gateway" "outscale_lin_internet_gateway" {}

resource "outscale_lin_internet_gateway_link" "outscale_lin_internet_gateway_link" {
  vpc_id              = "${outscale_lin.outscale_lin.id}"
  internet_gateway_id = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.id}"
}

# resource "outscale_nat_service" "outscale_nat_service" {
#   depends_on    = ["outscale_route.outscale_route"]
#   subnet_id     = "${outscale_subnet.outscale_subnet.subnet_id}"
#   allocation_id = "${outscale_public_ip.outscale_public_ip.allocation_id}"
# }


# resource "outscale_route_table" "outscale_route_table" {
#   vpc_id = "${outscale_lin.outscale_lin.id}"
# }


# resource "outscale_route" "outscale_route" {
#   destination_cidr_block = "0.0.0.0/0"
#   gateway_id             = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.id}"
#   route_table_id         = "${outscale_route_table.outscale_route_table.id}"
# }


# resource "outscale_route_table_link" "outscale_route_table_link" {
#   subnet_id      = "${outscale_subnet.outscale_subnet.subnet_id}"
#   route_table_id = "${outscale_route_table.outscale_route_table.id}"
# }


# resource "outscale_public_ip" "outscale_public_ip" {
#   domain = "vpc"
# }


# resource "outscale_public_ip_link" "by_public_ip" {
#   public_ip   = "${outscale_public_ip.outscale_public_ip.public_ip}"
#   instance_id = "${outscale_vm.basic.id}"
#   depends_on  = ["outscale_vm.basic", "outscale_public_ip.outscale_public_ip"]
# }

