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

data "outscale_public_ip" "outscale_public_ip_2" {
   filter {
      name  = "tags"
      values = ["name-1=public_ip-data"]
   }
}

data "outscale_public_ip" "outscale_public_ip_3" {
   filter {
      name  = "tag_keys"
      values = ["name-1"]
   }
}

data "outscale_public_ip" "outscale_public_ip_4" {
   filter {
      name  = "tag_values"
      values = ["public_ip-data"]
   }
}
