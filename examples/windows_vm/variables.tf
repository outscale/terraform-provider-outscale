variable "access_key_id" {}
variable "secret_key_id" {}
variable "region" {}

variable "image_id" {}
variable "vm_type" {}
variable "allowed_cidr" {
  type = list(string)
}
