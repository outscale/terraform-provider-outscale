resource "outscale_security_group" "my_sgImgs" {
  description         = "test sg-group"
  security_group_name = "test-sg-${random_string.suffix[0].result}"
}

resource "outscale_vm" "outscale_vm" {
  image_id           = var.image_id
  vm_type            = var.vm_type
  security_group_ids = [outscale_security_group.my_sgImgs.security_group_id]
}

resource "outscale_image" "outscale_image1" {
  image_name = "test-image-${random_string.suffix[0].result}"
  description = "TF-69"
  vm_id       = outscale_vm.outscale_vm.vm_id
  no_reboot   = "true"
  tags {
    key   = "Key:TF69"
    value = "value-tags"
  }
  tags {
    key   = "Key-2"
    value = "value:TF69"
  }
}

resource "outscale_image" "outscale_image2" {
  image_name = "test-image-${random_string.suffix[1].result}"
  vm_id      = outscale_vm.outscale_vm.vm_id
  no_reboot  = "true"
  tags {
    key   = "Key:TF69"
    value = "value:TF69"
  }
  tags {
    key   = "Key-2"
    value = "value-tags-2"
  }
}


data "outscale_images" "outscale_images" {

  filter {
    name   = "image_ids"
    values = [outscale_image.outscale_image1.image_id, outscale_image.outscale_image2.image_id]
  }
  filter {
    name   = "descriptions"
    values = [outscale_image.outscale_image1.description]
  }
}

data "outscale_images" "outscale_images_2" {
  filter {
    name   = "image_names"
    values = [outscale_image.outscale_image1.image_name, outscale_image.outscale_image2.image_name]
  }
  depends_on = [outscale_image.outscale_image1, outscale_image.outscale_image2]

}

data "outscale_images" "outscale_images_4" {
  filter {
    name   = "states"
    values = [outscale_image.outscale_image1.state]
  }
  filter {
    name   = "tag_keys"
    values = ["Key:TF69"]
  }
  depends_on = [outscale_image.outscale_image1, outscale_image.outscale_image2]
}

data "outscale_images" "outscale_images_5" {
  filter {
    name   = "states"
    values = [outscale_image.outscale_image1.state]
  }
  filter {
    name   = "tag_values"
    values = ["value:TF69"]
  }
  depends_on = [outscale_image.outscale_image1, outscale_image.outscale_image2]
}

data "outscale_images" "outscale_images_6" {
  filter {
    name   = "states"
    values = [outscale_image.outscale_image1.state]
  }
  filter {
    name   = "tags"
    values = ["Key:TF69=value:TF69"]
  }
  depends_on = [outscale_image.outscale_image1, outscale_image.outscale_image2]
}

data "outscale_images" "outscale_images_7" {
  filter {
    name   = "boot_modes"
    values = outscale_image.outscale_image1.boot_modes
  }
  filter {
    name   = "tags"
    values = ["Key:TF69=value:TF69"]
  }
  depends_on = [outscale_image.outscale_image1, outscale_image.outscale_image2]
}
