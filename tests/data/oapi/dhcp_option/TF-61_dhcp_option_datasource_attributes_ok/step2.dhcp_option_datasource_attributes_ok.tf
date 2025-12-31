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
       name   = "ntp_servers"
       values = [outscale_dhcp_option.outscale_dhcp_option.ntp_servers.0, outscale_dhcp_option.outscale_dhcp_option.ntp_servers.1]
    }
filter {
       name   = "domain_names"
       values = [outscale_dhcp_option.outscale_dhcp_option.domain_name]
    }
filter {
       name   = "domain_name_servers"
       values = [outscale_dhcp_option.outscale_dhcp_option.domain_name_servers.0, outscale_dhcp_option.outscale_dhcp_option.domain_name_servers.1]
    }
filter {
       name   = "tags"
       values = ["Name-1=test-MZI-3"]
    }
filter {
       name   = "tag_keys"
       values = ["Name-1"]
    }
filter {
       name   = "tag_values"
       values = ["test-MZI-3"]
    }
}

data "outscale_dhcp_option" "outscale_dhcp_option-2" {
filter {
       name   = "log_servers"
       values = [outscale_dhcp_option.outscale_dhcp_option.log_servers.0, outscale_dhcp_option.outscale_dhcp_option.log_servers.1]
    }
}
