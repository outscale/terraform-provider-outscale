variable "access_key_id" {}
variable "secret_key_id" {}
variable "region" {}

variable "allowed_cidr" {
  type = list(string)
}
