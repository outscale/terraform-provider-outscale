# provider configuration
variable "account_id" {}

# resources configuration
variable "image_id" {}
variable "region" {}

variable "vm_type" {
  type    = string
  default = "tinav4.c2r2p2"
}
variable "osu_bucket_name" {}
variable "service_name" {
  type    = string
  default = "com.outscale.eu-west-2.api"
}

variable "fgpu_gen" {
  type    = string
  default = "v3"
}

variable "fgpu_vm_type" {
  type    = string
  default = "tinav3.c2r2p2"
}
