resource "outscale_net" "outscale_net" {
    ip_range = "10.5.0.0/16"
}

resource "outscale_subnet" "subnet01" {
    subregion_name = "${var.region}a"
    ip_range       = "10.5.0.0/24"
    net_id         = outscale_net.outscale_net.net_id
}

resource "outscale_subnet" "subnet02" {
    subregion_name = "${var.region}b"
    ip_range       = "10.5.2.0/24"
    net_id         = outscale_net.outscale_net.net_id
}

resource "outscale_security_group" "security_group01" {
  description         = "Terraform security group test"
  security_group_name = "Terraform-SG"
  net_id              = outscale_net.outscale_net.net_id
}

resource "outscale_vm" "vm01" {
  image_id           = var.image_id
  vm_type            = var.vm_type
  subnet_id          = outscale_subnet.subnet01.subnet_id
}


resource "outscale_nic" "outscale_nic" {
    subnet_id          = outscale_subnet.subnet01.subnet_id
    description        = "TF-109"
    security_group_ids = [outscale_security_group.security_group01.security_group_id]
    private_ips {
      is_primary       = true
      private_ip       = "10.5.0.45"
    }

    private_ips {
      is_primary       = false
      private_ip       = "10.5.0.46"
    }
    tags {
      key              = "Key:"
      value            = ":value-tags"
    }
    tags {
      key              = "Key-2"
      value            = "value-tags-2"
    }
}

resource "outscale_nic" "outscale_nic_2" {
    subnet_id = outscale_subnet.subnet01.subnet_id
    private_ips {
      is_primary = true
      private_ip = "10.5.0.41"
    }
    private_ips { 
      is_primary = false
      private_ip = "10.5.0.42"
    }
    tags {             
      key        = "Name"
      value      = "Nic-2"
    }
}

resource "outscale_nic" "outscale_nic_3" {
    subnet_id = outscale_subnet.subnet02.subnet_id
   security_group_ids = [outscale_security_group.security_group01.security_group_id] 
   private_ips {
      is_primary = true
      private_ip = "10.5.2.21"
    }
    tags {
      key              = "Key:"
      value            = ":value-tags"
    }
    tags {
      key              = "Key-2"
      value            = "value-tags-2"
    }
}


resource "outscale_nic_link" "nic_link01" {
    device_number = "1"
    vm_id         = outscale_vm.vm01.vm_id
    nic_id        = outscale_nic.outscale_nic.nic_id
}


resource "outscale_nic_link" "nic_link02" {
    device_number = "2"
    vm_id         = outscale_vm.vm01.vm_id
    nic_id        = outscale_nic.outscale_nic_2.nic_id
}

data "outscale_nics" "nic-0" {
    filter {
        name = "nic_ids"
        values = [outscale_nic.outscale_nic.nic_id]
    }    
depends_on =[outscale_nic_link.nic_link01,outscale_nic_link.nic_link02]
}

data "outscale_nics" "nic-2-main" {
    filter {
        name = "link_nic_vm_ids"
        values = [outscale_nic_link.nic_link01.vm_id]
    }
   filter {
        name = "link_nic_device_numbers"
        values = [0]
    }
depends_on =[outscale_nic_link.nic_link01,outscale_nic_link.nic_link02]
}

data "outscale_nics" "nic-4" {
    filter {
        name = "private_ips_private_ips"
        values = ["10.5.0.42","10.5.2.21"]
    }
depends_on=[outscale_nic.outscale_nic,outscale_nic.outscale_nic_2,outscale_nic.outscale_nic_3,outscale_nic_link.nic_link01,outscale_nic_link.nic_link02]
}

data "outscale_nics" "nic-5" {
    filter {
        name = "security_group_ids"
        values = [outscale_security_group.security_group01.security_group_id]
    }
depends_on=[outscale_nic.outscale_nic,outscale_nic.outscale_nic_2,outscale_nic.outscale_nic_3,outscale_nic_link.nic_link01,outscale_nic_link.nic_link02]
}

data "outscale_nics" "nic-6" {
    filter {
        name = "security_group_names"
        values = [outscale_security_group.security_group01.security_group_name]
    }
depends_on=[outscale_nic.outscale_nic,outscale_nic.outscale_nic_2,outscale_nic.outscale_nic_3,outscale_nic_link.nic_link01,outscale_nic_link.nic_link02]
}

data "outscale_nics" "nic-7" {
    filter {
        name   = "subnet_ids"
        values = [outscale_nic.outscale_nic_3.subnet_id]
    }
depends_on=[outscale_nic.outscale_nic_3]
}

