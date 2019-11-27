resource "ouscale_lin" "outscale_lin" {
    count = 1

    cidr_block = "10.0.0.0/16"
}

resource "outscale_lin_attribute" "outscale_lin_attribute" {
    enable_dns_hostnames = True
    vpc_id = outscale_lin.outscale_lin.vpc_id
}
