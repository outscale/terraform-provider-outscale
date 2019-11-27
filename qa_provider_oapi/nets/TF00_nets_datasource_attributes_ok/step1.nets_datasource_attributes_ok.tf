resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/17"

    tags = {
        key   = "Name"
        value = "outscale_net_resource"
    }

    #provisioner "local-exec" {
    #    command = "ssh -o StrictHostKeyChecking=no root@${osc_eip.eip_instance.*.public_ip[count.index]} uname -a"
    #}
}

resource "outscale_net" "outscale_net2" {
    ip_range = "10.2.0.0/17"

    tags = {
        key   = "Name"
        value = "outscale_net_resource2"
    }
}

data "outscale_nets" "outscale_nets" {
    filter {
        name   = "net_ids"
        values = [outscale_net.outscale_net.net_id, outscale_net.outscale_net2.net_id]
    }
}
