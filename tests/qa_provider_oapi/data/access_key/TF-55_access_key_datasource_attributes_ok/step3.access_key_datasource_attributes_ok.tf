resource "outscale_access_key" "my_access_key"{
 state                  = "INACTIVE"
 }

data "outscale_access_key" "my_access_key"{
filter {
 name ="states"
 values = ["INACTIVE"]
  }
}
