variable "access_key_id" {}
variable "secret_key_id" {}
variable "region" {}

variable "volume_type" {}
variable "volume_iops" {}
variable "volume_size_gib" {}
variable "image_id" {}
variable "vm_type" {}
variable "allowed_cidr" {
  type = list(string)
}
