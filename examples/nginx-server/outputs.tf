output "public_ip" {
  description = "Public IP address of the Nginx server."
  value       = outscale_public_ip.my-public-ip.public_ip
}

output "web_url" {
  description = "HTTP URL of the deployed Nginx server."
  value       = "http://${outscale_public_ip.my-public-ip.public_ip}"
}

output "instance_name" {
  description = "Instance name used by the example."
  value       = var.instance_name
}

output "private_key_file" {
  description = "Path to the generated private SSH key."
  value       = local_file.private_key.filename
}

output "ssh_command" {
  description = "Suggested SSH command to connect to the instance."
  value       = "ssh -i ${local_file.private_key.filename} outscale@${outscale_public_ip.my-public-ip.public_ip}"
}