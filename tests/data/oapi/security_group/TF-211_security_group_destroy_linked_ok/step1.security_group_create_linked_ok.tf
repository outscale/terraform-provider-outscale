resource "outscale_net" "net" {
  ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet" {
  subregion_name = "eu-west-2a"
  ip_range       = "10.0.0.0/18"
  net_id         = outscale_net.net.net_id
}

resource "outscale_security_group" "group1" {
  security_group_name = "test-sg-${random_string.suffix[0].result}"
  description = "first security group linked resources"
  net_id = outscale_net.net.id
}

resource "outscale_security_group" "group2" {
  security_group_name = "test-sg-${random_string.suffix[1].result}"
  description = "second security group linked resources"
  net_id = outscale_net.net.id
}

resource "outscale_nic" "nic" {
  subnet_id = outscale_subnet.subnet.subnet_id
  security_group_ids = [outscale_security_group.group1.security_group_id, outscale_security_group.group2.security_group_id]
}

resource "outscale_vm" "vm" {
  image_id           = var.image_id
  vm_type            = var.vm_type
  keypair_name       = var.keypair_name
  security_group_ids = [outscale_security_group.group1.id, outscale_security_group.group2.id]
  subnet_id          = outscale_subnet.subnet.subnet_id
}

resource "outscale_load_balancer" "load_balancer02" {
  load_balancer_name = "test-lb-${random_string.suffix[0].result}"
  listeners {
    backend_port           = 80
    backend_protocol       = "TCP"
    load_balancer_protocol = "TCP"
    load_balancer_port     = 80
  }
  subnets            = [outscale_subnet.subnet.subnet_id]
  security_groups    = [outscale_security_group.group1.security_group_id, outscale_security_group.group2.security_group_id]
  load_balancer_type = "internal"
}
