resource "outscale_net" "outscale_net" {
    #count = 1

    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    #count = 1

    subregion_name = format("%s%s", var.region, "a")
    ip_range       = "10.0.0.0/16"
    net_id         = outscale_net.outscale_net.net_id
    tags {
    key = "name"
    value = "terraform-subnet"
    }
}

data "outscale_subnet" "outscale_subnet" {
    filter {
        name   = "subnet_ids"
        values = [outscale_subnet.outscale_subnet.subnet_id]
    }
}

#resource "outscale_security_group" "outscale_security_group" {
#description = "test Private VM"
#security_group_name = "Private-sg-group"
#net_id = outscale_net.outscale_net.net_id
#}

#resource "outscale_vm" "outscale_vm" {
#  image_id            = var.image_id
#  subnet_id =outscale_subnet.outscale_subnet.subnet_id
#  security_group_ids = [outscale_security_group.outscale_security_group.security_group_id]
#}
