variable "access_key_id" {}
variable "secret_key_id" {}
variable "region" {
    default = "eu-west-2"
}
variable "access_key" {}
variable "secret_key" {}
#variable "region" {
#    default = "eu-west-2"
#}

variable "ssh_privkey_filename" {
    description = "The name of the file containning the private ssh key"
    default     = "outscale_sutKeyPair.rsa.txt"
}

variable "vm_id" {}
variable "vm_type" {}
variable "image_id" {}
variable "volume_id" {}
variable "keypair_name" {}
variable "security_group_id" {}
variable "security_group_name" {}
variable "account_id" {}
variable "dhcp_options_set_id" {}
variable "snapshot_id" {}
