resource "outscale_virtual_gateway" "outscale_virtual_gateway" { 
 connection_type = "ipsec.1"  
 tags {   
  key = "name"   
  value = "test-VGW-1"   
  }
tags {
  key = "Project"
  value = "terraform"
  }
} 

