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

data "outscale_vpn_connection" "outscale_vpn_connection1" {
  vpn_connection_id = "${outscale_vpn_connection.outscale_vpn_connection.id}"
}

data "outscale_vpn_connections" "outscale_vpn_connections" {
  vpn_connection_id = ["${outscale_vpn_connection.outscale_vpn_connection.id}", "${outscale_vpn_connection.outscale_vpn_connection2.id}"]
}

resource "outscale_vpn_connection" "outscale_vpn_connection" {
  customer_gateway_id = "${outscale_client_endpoint.outscale_client_endpoint.id}"
  vpn_gateway_id      = "${outscale_vpn_gateway.outscale_vpn_gateway.id}"
  type                = "ipsec.1"
}

resource "outscale_vpn_gateway" "outscale_vpn_gateway" {
  type = "ipsec.1"
}

resource "outscale_client_endpoint" "outscale_client_endpoint" {
  bgp_asn    = "3"
  ip_address = "171.33.74.125"
  type       = "ipsec.1"
}

resource "outscale_vpn_gateway_link" "outscale_vpn_gateway_link" {
  vpc_id = "${outscale_lin.outscale_lin.vpc_id}"

  #vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.vpn_gateway_id}"
  vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.id}"
}

resource "outscale_lin" "outscale_lin" {
  cidr_block = "10.0.0.0/16"
}

resource "outscale_vpn_connection" "outscale_vpn_connection2" {
  customer_gateway_id = "${outscale_client_endpoint.outscale_client_endpoint2.id}"
  vpn_gateway_id      = "${outscale_vpn_gateway.outscale_vpn_gateway2.id}"
  type                = "ipsec.1"
}

resource "outscale_vpn_gateway" "outscale_vpn_gateway2" {
  type = "ipsec.1"
}

resource "outscale_client_endpoint" "outscale_client_endpoint2" {
  bgp_asn    = "3"
  ip_address = "171.33.74.126"
  type       = "ipsec.1"
}

resource "outscale_vpn_gateway_link" "outscale_vpn_gateway_link2" {
  vpc_id = "${outscale_lin.outscale_lin2.vpc_id}"

  #vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway2.vpn_gateway_id}"
  vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway2.id}"
}

resource "outscale_lin" "outscale_lin2" {
  cidr_block = "10.0.0.0/16"
}
