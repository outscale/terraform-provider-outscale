resource "outscale_access_key" "my_access_key"{
 state                  = "INACTIVE"
 }

data "outscale_access_key" "my_access_key"{
filter {
 name ="access_key_ids"
 values = [outscale_access_key.my_access_key.access_key_id]
  }
}
