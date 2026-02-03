# resources configuration
variable "image_id" {}
variable "region" {}
variable "osu_bucket_name" {}
variable "account_id" {}

variable "vm_type" {
  type    = string
  default = "tinav7.c2r2p1"
}

variable "fgpu_gen" {
  type    = string
  default = "v5"
}

variable "fgpu_vm_type" {
  type    = string
  default = "tinav5.c2r2p1"
}
