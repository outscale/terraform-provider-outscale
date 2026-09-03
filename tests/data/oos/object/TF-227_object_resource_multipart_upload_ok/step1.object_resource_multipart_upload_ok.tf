locals {
  file_name = "file.bin"
  path      = abspath("${path.module}/${local.file_name}")
}

resource "terraform_data" "file" {
  input = {
    path = local.path
  }

  provisioner "local-exec" {
    command = "dd if=/dev/zero of='${local.path}' bs=1M count=20 status=progress"
  }

  provisioner "local-exec" {
    when    = destroy
    command = "rm -f '${self.input.path}'"
  }
}

resource "outscale_oos_bucket" "multipart_test" {
  name = "test-oos-object-${random_string.suffix[0].result}"
}

resource "outscale_oos_object" "multipart_test" {
  bucket       = outscale_oos_bucket.multipart_test.id
  key          = local.file_name
  source       = local.file_name
  content_type = "application/octet-stream"

  depends_on = [terraform_data.file]
}
