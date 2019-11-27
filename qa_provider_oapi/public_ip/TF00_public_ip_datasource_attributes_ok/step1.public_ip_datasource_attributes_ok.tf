resource "outscale_public_ip" "outscale_public_ip" {
}

data "outscale_public_ip" "outscale_public_ip" {
   filter {
      name  = "public_ips"
      values = [outscale_public_ip.outscale_public_ip.public_ip]
   }    
}
