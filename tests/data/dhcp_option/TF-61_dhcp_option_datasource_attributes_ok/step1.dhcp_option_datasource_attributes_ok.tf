resource "outscale_dhcp_option" "outscale_dhcp_option" {
domain_name ="test234.fr"
domain_name_servers= ["192.168.12.12","192.168.12.132"]
ntp_servers = ["192.0.0.2","192.168.12.242"]
log_servers = ["192.10.10.2","192.168.112.92"]
tags {
   key ="Name-1"
   value = "test-MZI-3"
 }
tags {
   key ="Project"
   value = "terraform"
 }
}

data "outscale_dhcp_option" "outscale_dhcp_option" {
filter {
       name   = "dhcp_options_set_ids"
       values = [outscale_dhcp_option.outscale_dhcp_option.dhcp_options_set_id]
    }
}
