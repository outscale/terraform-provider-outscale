resource "outscale_dhcp_option" "outscale_dhcp_option" {
domain_name ="test-terraform.fr"
domain_name_servers= ["192.168.12.12","192.168.12.132"]
ntp_servers = ["192.0.0.2","192.168.12.242"]
log_servers = ["192.33.0.2","192.33.12.242"]
tags {
   key ="name-1"
   value = "test-terraform-3"
 }
tags {
   key ="Project"
   value = "test-terraform-3"
 }
}

resource "outscale_dhcp_option" "outscale_dhcp_option2" {
domain_name_servers = ["OutscaleProvidedDNS"]
tags {
   key ="name"
   value = "test-terraform_2"
 }
}

data "outscale_dhcp_options" "outscale_dhcp_options" {
filter {
       name   = "dhcp_options_set_ids"
       values = [outscale_dhcp_option.outscale_dhcp_option.dhcp_options_set_id]
    }
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
       values = ["name-1=test-terraform-3"]
    }
filter {
       name   = "tag_keys"
       values = ["name-1"]
    }
filter {
       name   = "tag_values"
       values = ["test-terraform-3"]
    }
}

data "outscale_dhcp_options" "outscale_dhcp_options_2" { 
filter {
       name   = "dhcp_options_set_ids"
       values = [outscale_dhcp_option.outscale_dhcp_option2.dhcp_options_set_id]
    }
}

data "outscale_dhcp_option" "outscale_dhcp_option-3" {
filter {
       name   = "log_servers"
       values = [outscale_dhcp_option.outscale_dhcp_option.log_servers.0, outscale_dhcp_option.outscale_dhcp_option.log_servers.1]
    }
}
