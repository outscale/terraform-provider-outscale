resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF144"
}

resource "outscale_security_group" "public_sg_terraform" {
    description         = "test vms_2"
    security_group_name = "terraform-public-sg_for_vms"
}
resource "outscale_vm" "outscale_vm" {
    image_id                 = var.image_id
    vm_type                  = var.vm_type
    keypair_name             = outscale_keypair.my_keypair.keypair_name
    security_group_ids       = [outscale_security_group.public_sg_terraform.security_group_id]
    placement_subregion_name = "${var.region}a"
    placement_tenancy        = "default"
    tags {
         key = "name"
         value = "outscale_vm"
         }
}

resource "outscale_vm" "outscale_vm2" {
    image_id                 = var.image_id
    vm_type                  = var.vm_type
    keypair_name             = outscale_keypair.my_keypair.keypair_name
    security_group_names     = [outscale_security_group.public_sg_terraform.security_group_name]
    placement_subregion_name = "${var.region}a"
    placement_tenancy        = "default"
     tags {
         key = "name"
         value = "outscale_vm2"
         }
}

# vm in net

resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"

    tags  {                               
        key   = "name"
        value = "Terraform_net"
      }
}

resource "outscale_subnet" "outscale_subnet" {
    net_id         = outscale_net.outscale_net.net_id
    ip_range       = "10.0.0.0/24"
    subregion_name = "${var.region}a"

    tags {                               
        key   = "name"
        value = "Terraform_subnet"
      }
}

resource "outscale_security_group" "outscale_sg" {
    description         = "sg for terraform tests"
    security_group_name = "terraform-sg"
    net_id              = outscale_net.outscale_net.net_id
     tags {                               
        key   = "name"
        value = "outscale_sg"
      }
} 
 
resource "outscale_internet_service" "outscale_internet_service" {
tags {                               
        key   = "name"
        value = "outscale_internet_service"
      }
depends_on = [outscale_net.outscale_net]
}

resource "outscale_route_table" "outscale_route_table" {
    net_id = outscale_net.outscale_net.net_id

    tags {                               
        key   = "name"
        value = "Terraform_RT"
      }
}

resource "outscale_route_table_link" "outscale_route_table_link" {
    route_table_id  = outscale_route_table.outscale_route_table.route_table_id
    subnet_id       = outscale_subnet.outscale_subnet.subnet_id
}

resource "outscale_internet_service_link" "outscale_internet_service_link" {
    internet_service_id = outscale_internet_service.outscale_internet_service.internet_service_id
    net_id              = outscale_net.outscale_net.net_id
}

resource "outscale_route" "outscale_route" {
    gateway_id           = outscale_internet_service.outscale_internet_service.internet_service_id
     destination_ip_range = "0.0.0.0/0"
    route_table_id       = outscale_route_table.outscale_route_table.route_table_id
} 

resource "outscale_vm" "outscale_vmnet" {
    image_id           = var.image_id
    vm_type            = "t2.nano"
    keypair_name       = outscale_keypair.my_keypair.keypair_name
    security_group_ids = [outscale_security_group.outscale_sg.security_group_id]
    subnet_id          = outscale_subnet.outscale_subnet.subnet_id
}

resource "outscale_public_ip" "outscale_public_ip" {
tags {                               
        key   = "name"
        value = "outscale_public_ip"
      }
}

resource "outscale_public_ip_link" "outscale_public_ip_link" {
    vm_id     = outscale_vm.outscale_vmnet.vm_id
    public_ip = outscale_public_ip.outscale_public_ip.public_ip
}
