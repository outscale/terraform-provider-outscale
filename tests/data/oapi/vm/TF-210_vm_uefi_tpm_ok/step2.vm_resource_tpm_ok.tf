resource "outscale_image" "image03" {
  image_name         = "test-image-${random_string.suffix[0].result}"
  source_image_id    = var.image_id_uefi_tpm
  source_region_name = var.region
  tpm_mandatory      = true
}

resource "outscale_vm" "outscale_vm_TF206" {
  image_id            = outscale_image.image03.id
  vm_type             = var.vm_type
  deletion_protection = false
  boot_mode           = "uefi"
  tpm_enabled = true
}
