# Generate an SSH keypair for the example.
# Note: the private key is stored in Terraform state.
resource "tls_private_key" "my-private-key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

# Register the generated public key in OUTSCALE.
resource "outscale_keypair" "my-keypair" {
  keypair_name = "${local.name_prefix}-keypair"
  public_key   = tls_private_key.my-private-key.public_key_openssh
}

# Save the generated private key locally so the instance can be accessed with SSH.
resource "local_file" "private_key" {
  filename        = "${path.module}/${var.private_key_filename}"
  content         = tls_private_key.my-private-key.private_key_pem
  file_permission = "0600"
}