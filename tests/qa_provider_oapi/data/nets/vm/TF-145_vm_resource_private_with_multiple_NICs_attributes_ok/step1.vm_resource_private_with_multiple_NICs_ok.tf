resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF145"
}
## Test Private VM with multiple NICs  ##

resource "outscale_net" "outscale_net" {
    ip_range = "10.22.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
 net_id              = outscale_net.outscale_net.net_id
 ip_range            = "10.22.0.0/24"
 subregion_name      = "${var.region}a"
}

resource "outscale_security_group" "outscale_security_group" {
    description         = "test vm with nic"
    security_group_name = "private-sg-1"
    net_id              = outscale_net.outscale_net.net_id
}

resource "outscale_security_group" "outscale_security_group2" {
    description         = "test vm with nic"
    security_group_name = "private-sg-2"
    net_id              = outscale_net.outscale_net.net_id
}


resource "outscale_nic" "outscale_nic" {
    subnet_id = outscale_subnet.outscale_subnet.subnet_id
}

resource "outscale_vm" "outscale_vm" {
    image_id            = var.image_id
    vm_type             = var.vm_type
    keypair_name        = outscale_keypair.my_keypair.keypair_name
    nics       {
       subnet_id = outscale_subnet.outscale_subnet.subnet_id
       security_group_ids = [outscale_security_group.outscale_security_group.security_group_id]
       private_ips  {
             private_ip ="10.22.0.123"
             is_primary = true
        }
       device_number = "0"
       delete_on_vm_deletion = true
    }
    nics {
       nic_id =outscale_nic.outscale_nic.nic_id
       device_number = "1"
    }
    tags {
       key = "name"
       value = "test-VM-with-Nics"   
    } 
 }

