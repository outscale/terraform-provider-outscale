resource "outscale_public_ip" "outscale_public_ip" {
tags {
 key = "name-1"
 value = "public_ip-data"
     }
tags {
      key = "Key"
      value = "value-tags"
     }
}

data "outscale_public_ip" "outscale_public_ip" {
   filter {
      name  = "public_ips"
      values = [outscale_public_ip.outscale_public_ip.public_ip]
   }    
}
