resource "outscale_image" "image03" {
  image_name         = "test-image-${random_string.suffix[0].result}"
  source_image_id    = var.image_id_uefi_tpm
  source_region_name = var.region
  tpm_mandatory      = true
}
