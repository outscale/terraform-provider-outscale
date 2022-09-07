resource "shell_script" "ca_gen" {
  lifecycle_commands {
    create = <<-EOF
          openssl req -x509 -sha256 -nodes -newkey rsa:4096 -keyout example.com.key -days 2 -out example.com.pem -subj '/CN=domain.com'
    EOF
    read   = <<-EOF
        echo "{\"filename\":  \"example.com.pem\"}"
    EOF
    delete = "rm -f example.com.pem example.com.key"
  }
  working_directory = "${path.module}/."
}

resource "outscale_ca" "my_ca" {
  ca_pem      = file(shell_script.ca_gen.output.filename)
  description = var.description
}
