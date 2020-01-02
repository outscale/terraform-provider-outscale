data "outscale_image" "outscale_image" {
    filter {
        name   = "image_ids"
        values = [var.image_id]
    }
}
