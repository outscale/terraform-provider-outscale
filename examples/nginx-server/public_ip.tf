# Allocate a public IP for the example instance.
resource "outscale_public_ip" "my-public-ip" {}

# Associate the public IP with the VM.
resource "outscale_public_ip_link" "my-public-ip-link" {
  public_ip = outscale_public_ip.my-public-ip.public_ip
  vm_id     = outscale_vm.my-vm.vm_id
}