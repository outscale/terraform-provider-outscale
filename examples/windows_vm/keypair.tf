resource "tls_private_key" "my_key" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "local_file" "my_key" {
  filename        = "${path.module}/my_key.pem"
  content         = tls_private_key.my_key.private_key_pem
  file_permission = "0600"
}

resource "outscale_keypair" "my_keypair" {
  public_key = tls_private_key.my_key.public_key_openssh
}
