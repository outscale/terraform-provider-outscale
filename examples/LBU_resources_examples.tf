#---Public LBU------------------------------------------------------------------
resource "outscale_load_balancer" "public_lbu-01" {
   load_balancer_name         = "lbu-01"
   subregion_names            = ["${var.region}a"]
   listeners {
      backend_port            = 80
      backend_protocol        = "HTTP"
      load_balancer_protocol  = "HTTP"
      load_balancer_port      = 80
    }
   listeners {
      backend_port            = 1024
      backend_protocol        = "TCP"
      load_balancer_protocol  = "TCP"
      load_balancer_port      = 1024
    }
   tags {
      key                     = "name"
      value                   = "public_lbu-01"
   }
}
#--------------------------------------------------------------------------------

#---Private LBU------------------------------------------------------------------
resource "outscale_net" "Net-01" {
    ip_range                  = "10.0.0.0/16"
}

resource "outscale_route_table" "Route_Table-01" {
    net_id                    = outscale_net.Net-01.net_id
}

resource "outscale_security_group" "SG-1" {
    description               = "test group"
    security_group_name       = "sg1-terraform-lbu-test"
    net_id                    = outscale_net.Net-01.net_id
}

resource "outscale_security_group" "SG-2" {
    description               = "test group-2"
    security_group_name       = "sg2-terraform-lbu-test"
    net_id                    = outscale_net.Net-01.net_id
}
resource "outscale_subnet" "Subnet-1" {
  net_id                      = outscale_net.Net-01.net_id
  ip_range                    = "10.0.0.0/24"
}

resource "outscale_load_balancer" "internal_lbu" {
   load_balancer_name         = "lbu-internal"
   subnets                    = [outscale_subnet.Subnet-1.subnet_id]
   security_groups            = [outscale_security_group.SG-1.id, outscale_security_group.SG-2.id]
   load_balancer_type         = "internal"
   listeners {
      backend_port            = 80
      backend_protocol        = "HTTP"
      load_balancer_protocol  = "HTTP"
      load_balancer_port      = 80
    }
   tags {
      key                     = "name"
      value                   = "internal_lbu"
   }
}
#----------------------------------------------------------------------------

#---LBU Policies-------------------------------------------------------------
resource "outscale_load_balancer" "public_lbu-02" {
   load_balancer_name        = "lbu-02"
   subregion_names           = ["${var.region}a"]
   listeners {
      backend_port           = 80
      backend_protocol       = "HTTP"
      load_balancer_protocol = "HTTP"
      load_balancer_port     = 80
    }
    tags {
      key                  = "project"
      value                = "terraform"
   }
}

resource "outscale_load_balancer_policy" "policy-1" {
    load_balancer_name      = outscale_load_balancer.public_lbu-02.load_balancer_name
    policy_name             = "lbu-policy-1"
    policy_type             = "load_balancer"
}

resource "outscale_load_balancer_policy" "policy-2" {
    load_balancer_name      = outscale_load_balancer.public_lbu-02.load_balancer_name
    policy_name             = "lbu-policy-2"
    policy_type             = "app"
    cookie_name             = "Cookie-2"
}

#----------------------------------------------------------------------------

#---LBU Backend_vms---------------------------------------------------------
resource "outscale_vm" "backend_vm" {
    count                  = 2  
    image_id               = var.image_id
    vm_type                = var.vm_type
    keypair_name           = var.keypair_name
}

resource "outscale_load_balancer" "public_lbu-03" {
  load_balancer_name       = "lbu-03"
  subregion_names          = ["${var.region}a"]
  listeners {
    backend_port           = 80
    backend_protocol       = "HTTP"
    load_balancer_protocol = "HTTP"
    load_balancer_port     = 80
   }
}

resource "outscale_load_balancer_vms" "backend_vms" {
   load_balancer_name      = outscale_load_balancer.public_lbu-03.load_balancer_name
   backend_vm_ids          = [outscale_vm.backend_vm[0].vm_id,outscale_vm.backend_vm[1].vm_id]
}

data "outscale_load_balancer_vm_health" "backend_vms_health" {
    load_balancer_name     = outscale_load_balancer.public_lbu-03.load_balancer_name
    backend_vm_ids         = [outscale_vm.backend_vm[0].vm_id,outscale_vm.backend_vm[1].vm_id]
depends_on = [outscale_load_balancer_vms.backend_vms]
}

#----------------------------------------------------------------------------

#---LBU Update Attributes----------------------------------------------------

resource "outscale_load_balancer" "public_lbu-04" {
  load_balancer_name        = "lbu-04"
  subregion_names           = ["${var.region}a"]
  listeners {
     backend_port           = 80
      backend_protocol      = "HTTP"
     load_balancer_protocol = "HTTP"
     load_balancer_port     = 80
   }
}

resource "outscale_load_balancer_attributes" "attributes-health_check" {
   load_balancer_name       = outscale_load_balancer.public_lbu-04.load_balancer_name
    health_check  {
        healthy_threshold   = 10
        check_interval      = 30
        path                = "/index.html"
        port                = 80
        protocol            = "HTTP"
        timeout             = 5
        unhealthy_threshold = 5
    }
}

resource "outscale_load_balancer_attributes" "attributes-access_log" {
   load_balancer_name = outscale_load_balancer.public_lbu-04.load_balancer_name
   access_log {
       publication_interval = 5
       is_enabled           = true
       osu_bucket_name      = "terraform-bucket"
       osu_bucket_prefix    = "access-logs"
   }
}

resource "outscale_load_balancer_attributes" "attributes-listener-policy" {
   load_balancer_name       = outscale_load_balancer.public_lbu-04.load_balancer_name
   load_balancer_port       = 80
   policy_names             = ["lbu-policy-3"]
depends_on =[outscale_load_balancer_policy.policy-3]
}
resource "outscale_load_balancer_policy" "policy-3" {
    load_balancer_name      = outscale_load_balancer.public_lbu-04.load_balancer_name
    policy_name             = "lbu-policy-3"
    policy_type             = "load_balancer"
}

#----------------------------------------------------------------------------

#---LBU datasource(s)--------------------------------------------------------

data "outscale_load_balancer" "public_lbu-01" {
  filter {
     name                   = "load_balancer_names"
     values                 = ["lbu-01"]
    }
depends_on = [outscale_load_balancer.public_lbu-01]
}

data "outscale_load_balancers" "public_lbus" {
  filter {
     name                   = "load_balancer_names"
     values                 = ["lbu-01","lbu-02"]
    }
depends_on = [outscale_load_balancer.public_lbu-01,outscale_load_balancer.public_lbu-02]
}

data "outscale_load_balancer_tags" "lbu_tags" {
  load_balancer_names      = ["lbu-01","lbu-02"]
depends_on = [outscale_load_balancer.public_lbu-01,outscale_load_balancer.public_lbu-02]
}

#----------------------------------------------------------------------------

#---Listener Rules-----------------------------------------------------------
resource "outscale_load_balancer_listener_rule" "rule-1" {
  listener {
     load_balancer_name     = outscale_load_balancer.public_lbu-03.load_balancer_name
     load_balancer_port     = 80
    }

  listener_rule {
     action                 = "forward"
     listener_rule_name     = "listener-rule-1"
     path_pattern           = "*.abc.*.abc.*.com"
     priority               = 10
    }
   vm_ids                   = [outscale_vm.backend_vm[0].vm_id]
}

resource "outscale_load_balancer_listener_rule" "rule-2" {
    listener  {
       load_balancer_name   = outscale_load_balancer.public_lbu-03.load_balancer_name
       load_balancer_port   = 80
    }

    listener_rule {
      action                = "forward"
      listener_rule_name    = "listener-rule-2"
      host_name_pattern     = "*.abc.-.abc.*.com"
      priority              = 1
    }
   vm_ids                   = [outscale_vm.backend_vm[1].vm_id]
}

data "outscale_load_balancer_listener_rule" "listener_rule" {
filter {
        name                = "listener_rule_names"
        values              = ["listener-rule-2"]
    }
depends_on =[outscale_load_balancer_listener_rule.rule-2]
}

data "outscale_load_balancer_listener_rules" "listener_rules" {
filter {
        name                = "listener_rule_names"
        values              = ["listener-rule-2","listener-rule-1"]
    }
depends_on =[outscale_load_balancer_listener_rule.rule-2,outscale_load_balancer_listener_rule.rule-1]
}
#----------------------------------------------------------------------------
