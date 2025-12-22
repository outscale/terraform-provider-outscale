resource "outscale_keypair" "my_keypair" {
  keypair_name = "test-keypair-${random_string.suffix[0].result}"
}
resource "outscale_public_ip" "public_ip01" {
  tags {
    key   = "name"
    value = "EIP-TF01"
  }
}

resource "outscale_security_group" "security_groupTF01" {
  description         = "Terraform"
  security_group_name = "test-sg-${random_string.suffix[0].result}"
}

resource "outscale_vm" "vm01" {
  image_id           = var.image_id
  security_group_ids = [outscale_security_group.security_groupTF01.security_group_id]
  vm_type            = var.vm_type
  keypair_name       = outscale_keypair.my_keypair.keypair_name
  tags {
    key   = "name"
    value = "vm-TF01"
  }
}

resource "outscale_public_ip_link" "public_ip_link01" {
  vm_id     = outscale_vm.vm01.vm_id
  public_ip = outscale_public_ip.public_ip01.public_ip
}


resource "outscale_load_balancer" "public_lbu1" {
  load_balancer_name = "test-lb-${random_string.suffix[0].result}"
  subregion_names    = ["${var.region}a"]
  listeners {
    backend_port           = 80
    backend_protocol       = "TCP"
    load_balancer_protocol = "TCP"
    load_balancer_port     = 80
  }
  listeners {
    backend_port           = 8080
    backend_protocol       = "HTTP"
    load_balancer_protocol = "HTTP"
    load_balancer_port     = 8080
  }
  tags {
    key   = "name"
    value = "public_lbu1"
  }
}


resource "outscale_vm" "vm02" { ###Create VMs for next steps
  count              = 4
  image_id           = var.image_id
  vm_type            = var.vm_type
  security_group_ids = [outscale_security_group.security_groupTF01.security_group_id]
  keypair_name       = outscale_keypair.my_keypair.keypair_name
  tags {
    key   = "name"
    value = "vm-TF01-${count.index}"
  }
}

resource "outscale_load_balancer_vms" "outscale_load_balancer_vms02" {
  load_balancer_name = outscale_load_balancer.public_lbu1.load_balancer_name
  backend_ips        = [for _, vm in outscale_vm.vm02 : vm.public_ip]
}



data "outscale_load_balancer" "load_balancer_TF01" {
  filter {
    name   = "load_balancer_names"
    values = [outscale_load_balancer.public_lbu1.load_balancer_name]
  }
  depends_on = [outscale_load_balancer_vms.outscale_load_balancer_vms02]
}

data "outscale_load_balancers" "load_balancers_TF01" {
  filter {
    name   = "load_balancer_names"
    values = [outscale_load_balancer.public_lbu1.load_balancer_name]
  }
  depends_on = [outscale_load_balancer_vms.outscale_load_balancer_vms02]
}
