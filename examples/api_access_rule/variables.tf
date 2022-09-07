variable "access_key_id" {}
variable "secret_key_id" {}
variable "region" {}

variable "ip_ranges" {
  type = list(string)
}
variable "description" {}
