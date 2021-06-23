resource "outscale_load_balancer" "public_lbu1" {
  load_balancer_name ="lbu-TF-79"
  subregion_names= ["${var.region}a"]
  listeners {
     backend_port = 8080
     backend_protocol= "HTTPS"
     load_balancer_protocol= "HTTPS"
     load_balancer_port = 8080
     server_certificate_id = "arn:aws:iam::866464343301:server-certificate/Cockpit/Certificate-2"
    }
}
