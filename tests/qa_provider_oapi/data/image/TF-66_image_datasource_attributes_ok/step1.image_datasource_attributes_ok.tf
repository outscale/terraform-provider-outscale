resource "outscale_security_group" "My-SGTF66" {
  description         = "test security group"
  security_group_name = "SG-TF66"
}

resource "outscale_vm" "my-vm" {
  vm_type            = var.vm_type
  image_id           = var.image_id
  security_group_ids = [outscale_security_group.My-SGTF66.security_group_id]
}

resource "outscale_image" "outscale_image" {
  image_name  = "image-TF66"
  vm_id       = outscale_vm.my-vm.vm_id
  description = "TF-66"
  no_reboot   = "true"
  tags {
    key   = "Key:TF66"
    value = "value:TF66"
  }
  tags {
    key   = "Key-2"
    value = "value-tags-2"
  }
}


data "outscale_image" "outscale_image" {
  filter {
    name   = "image_ids"
    values = [outscale_image.outscale_image.image_id]
  }
}

data "outscale_image" "outscale_image_2" {
  filter {
    name   = "image_names"
    values = [outscale_image.outscale_image.image_name]
  }
  depends_on = [outscale_image.outscale_image]
}

data "outscale_image" "outscale_image_3" {
  filter {
    name   = "descriptions"
    values = [outscale_image.outscale_image.description]
  }
  depends_on = [outscale_image.outscale_image]
}

data "outscale_image" "outscale_image_4" {
  filter {
    name   = "states"
    values = [outscale_image.outscale_image.state]
  }
  filter {
    name   = "tag_keys"
    values = ["Key:TF66"]
  }
}

data "outscale_image" "outscale_image_5" {
  filter {
    name   = "states"
    values = [outscale_image.outscale_image.state]
  }
  filter {
    name   = "tag_values"
    values = ["value:TF66"]
  }
}

data "outscale_image" "outscale_image_6" {
  filter {
    name   = "states"
    values = [outscale_image.outscale_image.state]
  }
  filter {
    name   = "tags"
    values = ["Key:TF66=value:TF66"]
  }
}

data "outscale_image" "outscale_image_7" {
  filter {
    name   = "product_codes"
    values = [outscale_image.outscale_image.product_codes[0]]
  }
  filter {
    name   = "descriptions"
    values = [outscale_image.outscale_image.description]
  }
}
