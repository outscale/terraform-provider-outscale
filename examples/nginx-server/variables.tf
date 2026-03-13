variable "access_key" {
  description = "OUTSCALE API access key."
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "OUTSCALE API secret key."
  type        = string
  sensitive   = true
}

variable "region" {
  description = "Region where resources are created."
  type        = string
  default     = "eu-west-2"
}

variable "image_id" {
  description = "Image ID used to create the VM."
  type        = string
}

variable "vm_type" {
  description = "VM type used for the example."
  type        = string
  default     = "tinav5.c2r2p2"
}

variable "project_name" {
  description = "Project name displayed in tags and the demo web page."
  type        = string
  default     = "terraform-nginx-demo"
}

variable "instance_name" {
  description = "Instance name displayed in tags and the demo web page."
  type        = string
  default     = "nginx-server"
}

variable "root_volume_size" {
  description = "Root volume size in GiB."
  type        = number
  default     = 20
}

variable "private_key_filename" {
  description = "Local filename used to save the generated private SSH key."
  type        = string
  default     = "id_rsa_nginx_example"
}