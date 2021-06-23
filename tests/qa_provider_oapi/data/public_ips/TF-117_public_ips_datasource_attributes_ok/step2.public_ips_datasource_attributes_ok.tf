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

data "outscale_public_ips" "outscale_public_ips_2" {
   filter {
      name  = "tags"
      values = ["name-2=test-data-2"]
   }
}

data "outscale_public_ips" "outscale_public_ips_3" {
   filter {
      name  = "tag_keys"
      values = ["name-2"]
   }
}

data "outscale_public_ips" "outscale_public_ips_4" {
   filter {
      name  = "tag_values"
      values = ["test-data-2"]
   }
}
