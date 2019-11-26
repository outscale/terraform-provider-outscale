resource "outscale_net" "outscale_net" {
    count = 1

    ip_range = "10.0.0.0/16"
}

output "net" {
    value = "${outscale_net.outscale_net.net_id}"
}

