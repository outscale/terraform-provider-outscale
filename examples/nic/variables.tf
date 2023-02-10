variable "access_key_id" {}
variable "secret_key_id" {}
variable "region" {}

variable "image_id" {}
variable "vm_type" {}
variable "allowed_cidr" {
  type = list(string)
}
variable "net_ip_range" {}
variable "subnet_public_ip_range" {}
variable "subnet_private_ip_range" {}
variable "customer_ip_range" {}
