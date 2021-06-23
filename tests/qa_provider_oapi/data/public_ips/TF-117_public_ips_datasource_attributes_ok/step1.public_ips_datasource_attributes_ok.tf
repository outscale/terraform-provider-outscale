resource "outscale_public_ip" "outscale_public_ip" {
tags {
  key = "name-2"
  value = "test-data-2"
 }
tags {
      key = "Key"
      value = "value-tags"
     }
}

resource "outscale_public_ip" "outscale_public_ip2" {
}

data "outscale_public_ips" "outscale_public_ips" {
   filter {
      name  = "public_ips"
      values = [outscale_public_ip.outscale_public_ip.public_ip,outscale_public_ip.outscale_public_ip2.public_ip]
   }    
}
