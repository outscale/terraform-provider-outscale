resource "outscale_net" "outscale_net" {
     ip_range      = "10.0.0.0/16"
}

resource "outscale_net_access_point" "net_access_point_1" {
   net_id          = outscale_net.outscale_net.net_id
   service_name    = "com.outscale.${var.region}.api"
  tags {
          key      = "name"
          value    = "terraform-Net-Access-Point"
   }
   tags {
          key      = "test"
          value    = "test-Net-Access-Point"
   }
}
