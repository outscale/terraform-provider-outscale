data "outscale_vms" "vmd22" {                    # website
   filter {
      name = "image_id"                          # invalid filter
      values = ["ami-880caa66"]
   }
}

