resource "outscale_dhcp_option" "dhcp_option_1" {
domain_name ="test234.fr"
domain_name_servers= ["192.168.12.12","192.168.12.132"]
ntp_servers = ["192.0.0.2","192.168.12.242"]
tags {
   key ="name-1"
   value = "set-1"
 }
}

resource "outscale_dhcp_option" "dhcp_option_2" {
domain_name ="test2.fr"
tags {
   key ="name-2"
   value = "set-2"
 }
}


resource "outscale_dhcp_option" "dhcp_option_3" {
domain_name_servers= ["192.168.12.32","192.168.12.33"]
tags {
   key ="name-3"
   value = "set-3"
 }
}

resource "outscale_dhcp_option" "dhcp_option_4" {
ntp_servers = ["192.0.0.25","192.168.12.24"]
tags {
   key ="name-4"
   value = "set-4"
 }
}
