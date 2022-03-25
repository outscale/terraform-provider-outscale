resource "outscale_net" "outscale_net" {
    #count = 1

    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    subregion_name  = "${var.region}a"
    ip_range        = "10.0.0.0/16"
    net_id          = outscale_net.outscale_net.net_id
}

resource "outscale_nic" "outscale_nic" {
   
    subnet_id = outscale_subnet.outscale_subnet.subnet_id
    security_group_ids = [outscale_security_group.outscale_sg.security_group_id]
    private_ips {
     is_primary = true
     private_ip = "10.0.67.45"
   }
}

resource "outscale_security_group" "outscale_sg" {
    description         = "sg for terraform tests"
    security_group_name = "terraform-sg"
    net_id              = outscale_net.outscale_net.net_id
}

resource "outscale_nic" "outscale_nic2" {
    subnet_id = outscale_subnet.outscale_subnet.subnet_id
    security_group_ids = [outscale_security_group.outscale_sg.security_group_id]
    private_ips {
     is_primary = true
     private_ip = "10.0.0.23"
   }
     private_ips {
     is_primary = false
     private_ip = "10.0.0.46"
   }
}
