# provider configuration
variable "account_id" {}

# resources configuration
variable "image_id" {}
variable "region" {}
variable "service_name" {}
variable "osu_bucket_name" {}

variable "vm_type" {
  type    = string
  default = "tinav4.c2r2p2"
}

variable "fgpu_gen" {
  type    = string
  default = "v5"
}

variable "fgpu_vm_type" {
  type    = string
  default = "tinav5.c2r2p1"
}
